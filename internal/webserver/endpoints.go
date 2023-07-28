package webserver

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/Binozo/GoDrop/internal/awdl"
	"github.com/Binozo/GoDrop/internal/interaction"
	"github.com/Binozo/GoDrop/internal/utils"
	"github.com/Binozo/GoDrop/pkg/air"
	"github.com/korylprince/go-cpio-odc"
	"github.com/rs/zerolog/log"
	"howett.net/plist"
	"io"
	"net/http"
	"strings"
)

// HEAD /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Connection", "close")
	w.Header().Add("Content-Length", "0")
	w.WriteHeader(http.StatusOK)
}

// POST /Discover
func discoverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(awdl.DiscoverResponse)
}

// POST /Ask
func askHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Malformed request
		log.Warn().Err(err).Msg("Couldn't read 'POST /Ask' request")
		return
	}

	// Decoding the plist
	decoder := plist.NewDecoder(bytes.NewReader(body))
	var plistData map[string]interface{}
	decoder.Decode(&plistData)

	// Now we can gather some interesting data like the FileIcon or BundleID
	senderApplication := plistData["BundleID"].(string) // On a Mac it's 'com.apple.finder'
	//convertMediaFormats := plistData["ConvertMediaFormats"].(bool) // I don't know what this is used for
	fileIconRaw := plistData["FileIcon"].([]byte) // FileIcon is a preview of the to-be-sent file. Filetype is jp2
	// We will convert fileIconRaw later

	//items := plistData["Items"] // I don't know what this is used for
	senderComputerName := plistData["SenderComputerName"].(string)
	senderID := plistData["SenderID"].(string)
	senderModelName := plistData["SenderModelName"].(string)

	filesRaw := plistData["Files"].([]interface{})
	var files = make([]map[string]interface{}, 0)
	for _, fileRaw := range filesRaw {
		files = append(files, fileRaw.(map[string]interface{}))
	}

	// first we build the []File containing all "Files"
	airDropFiles := make([]air.File, 0)
	for _, fileData := range files {
		fileIsADirectory := fileData["FileIsDirectory"].(bool)
		fileName := fileData["FileName"].(string)
		fileBomPath := fileData["FileBomPath"].(string) // directory structure
		//fileType := files["FileType"]       // example: "public.png" or "public.folder" I don't know what "public" means in this context

		file := air.File{
			FileName:    fileName,
			FileBomPath: fileBomPath,
			IsDirectory: fileIsADirectory,
		}

		airDropFiles = append(airDropFiles, file)
	}

	fileIcon := []byte{}

	// Converting fileIconRaw (jp2) -> fileIcon (png)
	imageMagickInstalled := utils.IsImageMagickInstalled()
	if imageMagickInstalled {
		// yay it is installed :)
		// let's convert
		fileIcon, err = utils.ConvertJP2ToPNG(fileIconRaw)
		if err != nil {
			// Failed somehow to convert
			fileIcon = []byte{}
			log.Error().Err(err).Msg("Image conversion from jp2 to png failed")
		}
	} else {
		log.Warn().Msg("Couldn't convert preview file icon from jp2 to png: ImageMagick is not installed")
	}

	// Now we finished our decoding process
	airDropRequest := air.Request{
		SenderApplication:  senderApplication,
		SenderComputerName: senderComputerName,
		SenderID:           senderID,
		SenderModelName:    senderModelName,
		FileIcon:           fileIcon,
		Files:              airDropFiles,
		SenderIP:           parseIPv6(r),
	}

	accepted := interaction.AskUserAboutAcceptingIncomingFiles(airDropRequest)
	if accepted {
		// let's get that nice data
		w.Header().Add("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(awdl.DiscoverResponse)
	} else {
		// nah keep it for yourself
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// POST /Upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Now we check if the file ends with .gz
	var cpioArchive *bytes.Reader
	if data[0] == 0x1f && data[1] == 0x8b {
		// extract archive
		b := bytes.NewBuffer(data)
		var r io.Reader
		r, err := gzip.NewReader(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cpioArchive = bytes.NewReader(resB.Bytes())
	} else {
		cpioArchive = bytes.NewReader(data)
	}

	reader := cpio.NewReader(cpioArchive)
	var files []*cpio.File

	// Now we filter all unnecessary stuff out
	var file *cpio.File
	for file, err = reader.Next(); err == nil; file, err = reader.Next() {
		// We don't want files with names like "." or "._" ...
		if file.Name() != "." && !strings.HasPrefix(file.Name(), "._") {
			files = append(files, file)
		}
	}
	if errors.Is(err, io.EOF) {
		// Success
		interaction.OnFiles(files, parseIPv6(r))

	} else if err != nil {
		log.Warn().Err(err).Msg("Couldn't read cpio archive")
	}

	w.Header().Add("Connection", "close")
	w.WriteHeader(http.StatusOK)
}
