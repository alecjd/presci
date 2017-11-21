package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
	"gonode"
	"json"
	"github.com/gopherjs/gopherjs/js"
	
)

// Prefetch file header struct
type PrefetchHeader struct {
	Version   uint32
	Signature string
	_unk0     []byte
	Filesize  uint32
	Exename   string
	Hash      string
	_unk1     []byte
}

// Windows XP prefetech file information struct
type FileInfo17 struct {
	MetricsOffset         uint32
	MetricsCount          uint32
	TraceChainsOffset     uint32
	TraceChainsCount      uint32
	FilenameStringsOffset uint32
	FilenameStringsSize   uint32
	VolumesInfoOffset     uint32
	VolumesCount          uint32
	VolumesInfoSize       uint32
	LastRunTime           time.Time
	_unk2                 []byte
	RunCount              uint32
	_unk3                 uint32
}

// Windows 7 prefetch file information struct
type FileInfo23 struct {
	MetricsOffset         uint32
	MetricsCount          uint32
	TraceChainsOffset     uint32
	TraceChainsCount      uint32
	FilenameStringsOffset uint32
	FilenameStringsSize   uint32
	VolumesInfoOffset     uint32
	VolumesCount          uint32
	VolumesInfoSize       uint32
	_unk4                 []byte
	LastRunTime           time.Time
	_unk5                 []byte
	RunCount              uint32
	_unk6                 []byte
}

// Windows 8.1 prefetch file information struct
type FileInfo26 struct {
	MetricsOffset         uint32
	MetricsCount          uint32
	TraceChainsOffset     uint32
	TraceChainsCount      uint32
	FilenameStringsOffset uint32
	FilenameStringsSize   uint32
	VolumesInfoOffset     uint32
	VolumesCount          uint32
	VolumesInfoSize       uint32
	_unk7                 []byte
	LastRunTime           time.Time
	_unk8                 []byte
	RunCount              uint32
	_unk9                 []byte
}

type MetricsArray17 struct {
	_unk10         []byte
	_unk11         []byte
	FilenameOffset uint32
	FilenameLength uint32
	_unk12         []byte
}

type MetricsArray23 struct {
	_unk13            []byte
	_unk14            []byte
	_unk15            []byte
	FilenameOffset    uint32
	FilenameLength    uint32
	_unk16            []byte
	MFTRecordNumber   []byte
	MFTSequenceNumber []byte
}

type MetricsArray30 struct {
	_unk17         []byte
	_unk18         []byte
	FilenameOffset uint32
	FilenameLength uint32
	_unk19         []byte
}

type TraceChainsArray17 struct {
	_unk20 []byte
}

type TraceChainsArray30 struct {
	_unk21 []byte
}

type VolumeInfo17 struct {
	VolumePathOffset   uint32
	VolumePathLength   uint32
	VolumeCreationTime time.Time
	VolumeSerialNumber string
	FileRefOffset      uint32
	FileRefSize        uint32
	DirStringsOffset   uint32
	DirStringsCount    uint32
	_unk22             []byte
}

type VolumeInfo23 struct {
	VolumePathOffset   uint32
	VolumePathLength   uint32
	VolumeCreationTime time.Time
	VolumeSerialNumber string
	FileRefOffset      uint32
	FileRefCount       uint32
	DirStringsOffset   uint32
	DirStringsCount    uint32
	_unk23             []byte
}

type VolumeInfo30 struct {
	VolumePathOffset   uint32
	VolumePathLength   uint32
	VolumeCreationTime time.Time
	VolumeSerialNumber string
	FileRefOffset      uint32
	FileRefSize        uint32
	DirStringsOffset   uint32
	DirStringsCount    uint32
	_unk24             []byte
}

type PrefetchFileV17 struct {
	Header17    PrefetchHeader
	FileInfo17  FileInfo17
	MetArr17    MetricsArray17
	TChainArr17 TraceChainsArray17
	VolInfo17   []VolumeInfo17
	DirStrings  []string
	FnStrings   []string
}

const volinfo17_entry int = 40
const volinfo23_entry int = 104
const volinfo30_entry int = 96

func main() {
	js.Module.Get("exports").Set("pet", map[string]interface{}{
	  "New": ,
	}
	// var pfData []PfSection
	// path := "C:/Users/tailwindfor/test_data/"
	// get path from a config file or some sort of fixed location

	//this is all moved out into the function that main calls for the nodejs wrapper
	//the path needs to be more flexible to be selected from the node-based interface
	//ensure the node libraries are properly imported to ensure a proper connection
	path := "C:/Users/tailwindfor/test_data/XPPro"
	pf_file_paths := getPrefetchFiles(path)
	var pf_header = PrefetchHeader{}
	var pf_finfo17 = FileInfo17{}
	// var pf_finfo23 = FileInfo23{}
	// var pf_finfo26 = FileInfo26{}
	var pf_metarr17 = MetricsArray17{}
	// var pf_metarr23 = MetricsArray23{}
	// var pf_metarr30 = MetricsArray30{}
	var pf_tchainsarr17 = TraceChainsArray17{}
	// var pf_tchainsarr30 = TraceChainsArray30{}
	// var pf_volinfo17 = VolumeInfo17{}
	// var pf_volinfo23 = VolumeInfo23{}
	// var pf_volinfo30 = VolumeInfo30{}
	var pf_v17 = PrefetchFileV17{}

	//this code should be deferred to electron
	for _, pf := range pf_file_paths {
		// var vols = []VolumeInfo17{}
		pf_handle := openPrefetchFile(pf)
		pf_header.getPrefetchFileHeader(pf_handle)
		switch version := fmt.Sprint(pf_header.Version); version {
		case "17":
			pf_finfo17.getFileInfo17(pf_handle)
			pf_metarr17.getMetricsArray17(pf_handle)
			pf_tchainsarr17.getTraceChainsArray17(pf_handle)
			pf_v17.Header17 = pf_header
			pf_v17.FileInfo17 = pf_finfo17
			pf_v17.MetArr17 = pf_metarr17
			pf_v17.TChainArr17 = pf_tchainsarr17
			pf_v17.VolInfo17, pf_v17.DirStrings = getVolInfo17(pf_handle, pf_finfo17.VolumesCount, pf_finfo17.VolumesInfoOffset)
			pf_v17.FnStrings = getFileNameStrings(pf_handle, pf_finfo17.FilenameStringsSize, pf_finfo17.FilenameStringsOffset)
			j, err := json.MarshalIndent(&pf_v17, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(j))
			// fmt.Printf("%+v", vols)
			// fmt.Println("Windows XP/2003")
		case "23":
			fmt.Println("Windows Vista/7")
		case "26":
			fmt.Println("Windows8.1")
		case "30":
			fmt.Println("Windows 10")
		}
		closePrefetchFile(pf, pf_handle)
	}
}

func (hdr *PrefetchHeader) MarshalJSON2() []byte {
	j, err := json.MarshalIndent(&PrefetchHeader{}, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return j
}

//Prefetch file header information
//84 bytes
func (hdr *PrefetchHeader) getPrefetchFileHeader(pf_file *os.File) {
	data := readBytes(pf_file, 0, 84)
	hdr.Version = binary.LittleEndian.Uint32(data[:4])
	hdr.Signature = string(data[4:8])
	hdr._unk0 = data[8:12]
	hdr.Filesize = binary.LittleEndian.Uint32(data[12:16])
	hdr.Exename = formatExeName(data[16:76])
	hdr.Hash = hex.EncodeToString(data[76:80])
	hdr._unk1 = data[80:84]
}

//Windows XP (version: 17) prefetch file information
//68 bytes
func (finfo17 *FileInfo17) getFileInfo17(pf_file *os.File) {
	data := readBytes(pf_file, 84, 152)
	finfo17.MetricsOffset = binary.LittleEndian.Uint32(data[:4])
	finfo17.MetricsCount = binary.LittleEndian.Uint32(data[4:8])
	finfo17.TraceChainsOffset = binary.LittleEndian.Uint32(data[8:12])
	finfo17.TraceChainsCount = binary.LittleEndian.Uint32(data[12:16])
	finfo17.FilenameStringsOffset = binary.LittleEndian.Uint32(data[16:20])
	finfo17.FilenameStringsSize = binary.LittleEndian.Uint32(data[20:24])
	finfo17.VolumesInfoOffset = binary.LittleEndian.Uint32(data[24:28])
	finfo17.VolumesCount = binary.LittleEndian.Uint32(data[28:32])
	finfo17.VolumesInfoSize = binary.LittleEndian.Uint32(data[32:36])
	finfo17.LastRunTime = formatWin32FileTime(binary.LittleEndian.Uint64(data[36:44]))
	finfo17._unk2 = data[44:60]
	finfo17.RunCount = binary.LittleEndian.Uint32(data[60:64])
	finfo17._unk3 = binary.LittleEndian.Uint32(data[64:68])
}

//Windows XP (version: 17) prefetch metrics array information
//20 bytes
func (metarr17 *MetricsArray17) getMetricsArray17(pf_file *os.File) {
	data := readBytes(pf_file, 152, 172)
	metarr17._unk10 = data[:4]
	metarr17._unk11 = data[4:8]
	metarr17.FilenameOffset = binary.LittleEndian.Uint32(data[8:12])
	metarr17.FilenameLength = binary.LittleEndian.Uint32(data[12:16])
	metarr17._unk12 = data[16:20]
}

//Windows XP (version: 17) prefetch trace chains array information
//12 bytes
func (tchainsarr17 *TraceChainsArray17) getTraceChainsArray17(pf_file *os.File) {
	data := readBytes(pf_file, 172, 184)
	tchainsarr17._unk20 = data[:12]
}

//Windows 7 (version: 23) prefetch file information
//156 bytes
func (finfo23 *FileInfo23) getFileInfo23(pf_file *os.File) {
	data := readBytes(pf_file, 84, 240)
	finfo23.MetricsOffset = binary.LittleEndian.Uint32(data[:4])
	finfo23.MetricsCount = binary.LittleEndian.Uint32(data[4:8])
	finfo23.TraceChainsOffset = binary.LittleEndian.Uint32(data[8:12])
	finfo23.TraceChainsCount = binary.LittleEndian.Uint32(data[12:16])
	finfo23.FilenameStringsOffset = binary.LittleEndian.Uint32(data[16:20])
	finfo23.FilenameStringsSize = binary.LittleEndian.Uint32(data[20:24])
	finfo23.VolumesInfoOffset = binary.LittleEndian.Uint32(data[24:28])
	finfo23.VolumesCount = binary.LittleEndian.Uint32(data[28:32])
	finfo23.VolumesInfoSize = binary.LittleEndian.Uint32(data[32:36])
	finfo23._unk4 = data[36:44]
	finfo23.LastRunTime = formatWin32FileTime(binary.LittleEndian.Uint64(data[44:52]))
	finfo23._unk5 = data[52:68]
	finfo23.RunCount = binary.LittleEndian.Uint32(data[68:72])
	finfo23._unk6 = data[72:156]
}

//Windows 7 (version: 23) prefetch metrics array information
//32 bytes
func (metarr23 *MetricsArray23) getMetricsArray23(pf_file *os.File) {
	data := readBytes(pf_file, 240, 272)
	metarr23._unk13 = data[:4]
	metarr23._unk14 = data[4:8]
	metarr23._unk15 = data[8:12]
	metarr23.FilenameOffset = binary.LittleEndian.Uint32(data[12:16])
	metarr23.FilenameLength = binary.LittleEndian.Uint32(data[16:20])
	metarr23._unk16 = data[20:24]
	metarr23.MFTRecordNumber = data[24:30]   //fix
	metarr23.MFTSequenceNumber = data[30:32] //fix
}

//Windows 8.1 (version: 26) prefetch file information
//224 bytes
func (finfo26 *FileInfo26) getFileInfo26(pf_file *os.File) {
	data := readBytes(pf_file, 84, 308)
	finfo26.MetricsOffset = binary.LittleEndian.Uint32(data[:4])
	finfo26.MetricsCount = binary.LittleEndian.Uint32(data[4:8])
	finfo26.TraceChainsOffset = binary.LittleEndian.Uint32(data[8:12])
	finfo26.TraceChainsCount = binary.LittleEndian.Uint32(data[12:16])
	finfo26.FilenameStringsOffset = binary.LittleEndian.Uint32(data[16:20])
	finfo26.FilenameStringsSize = binary.LittleEndian.Uint32(data[20:24])
	finfo26.VolumesInfoOffset = binary.LittleEndian.Uint32(data[24:28])
	finfo26.VolumesCount = binary.LittleEndian.Uint32(data[28:32])
	finfo26.VolumesInfoSize = binary.LittleEndian.Uint32(data[32:36])
	finfo26._unk7 = data[36:44]
	finfo26.LastRunTime = formatWin32FileTime(binary.LittleEndian.Uint64(data[44:108]))
	finfo26._unk8 = data[108:124]
	finfo26.RunCount = binary.LittleEndian.Uint32(data[124:128])
	finfo26._unk9 = data[128:224]
}

//Windows XP (version: 17) prefetch volume information
//xx bytes
func getVolInfo17(pf_file *os.File, count uint32, volinfo_offset uint32) ([]VolumeInfo17, []string) {
	cnt := 0
	vol_info := []VolumeInfo17{}
	vol_data := VolumeInfo17{}
	dirstrings := []string{}
	end_offset := int(volinfo_offset) + int(count)*40

	for cnt < int(count) {
		data := readBytes(pf_file, int(volinfo_offset), end_offset)
		vol_data.VolumePathOffset = binary.LittleEndian.Uint32(data[:4])
		vol_data.VolumePathLength = binary.LittleEndian.Uint32(data[4:8])
		vol_data.VolumeCreationTime = formatWin32FileTime(binary.LittleEndian.Uint64(data[8:16]))
		vol_data.VolumeSerialNumber = hex.EncodeToString(data[16:20])
		vol_data.FileRefOffset = binary.LittleEndian.Uint32(data[20:24])
		vol_data.FileRefSize = binary.LittleEndian.Uint32(data[24:28])
		vol_data.DirStringsOffset = binary.LittleEndian.Uint32(data[28:32])
		vol_data.DirStringsCount = binary.LittleEndian.Uint32(data[32:36])
		vol_data._unk22 = data[36:40]
		cnt += 1
		vol_info = append(vol_info, vol_data)
		dirstrings = getDirectoryStrings(pf_file, volinfo_offset, vol_data.DirStringsCount, vol_data.DirStringsOffset)
	}
	return vol_info, dirstrings
}

//Windows 7 (version: 23) prefetch volume information
//xx bytes
func getVolInfo23(pf_file *os.File, count uint32, volinfo_offset uint32) ([]VolumeInfo23, []string) {
	cnt := 0
	vol_info := []VolumeInfo23{}
	vol_data := VolumeInfo23{}
	dirstrings := []string{}
	end_offset := int(volinfo_offset) + int(count)*104

	for cnt < int(count) {
		data := readBytes(pf_file, int(volinfo_offset), end_offset)
		vol_data.VolumePathOffset = binary.LittleEndian.Uint32(data[:4])
		vol_data.VolumePathLength = binary.LittleEndian.Uint32(data[4:8])
		vol_data.VolumeCreationTime = formatWin32FileTime(binary.LittleEndian.Uint64(data[8:16]))
		vol_data.VolumeSerialNumber = hex.EncodeToString(data[16:20])
		vol_data.FileRefOffset = binary.LittleEndian.Uint32(data[20:24])
		vol_data.FileRefCount = binary.LittleEndian.Uint32(data[24:28])
		vol_data.DirStringsOffset = binary.LittleEndian.Uint32(data[28:32])
		vol_data.DirStringsCount = binary.LittleEndian.Uint32(data[32:36])
		vol_data._unk23 = data[36:104]
		cnt += 1
		vol_info = append(vol_info, vol_data)
		dirstrings = getDirectoryStrings(pf_file, volinfo_offset, vol_data.DirStringsCount, vol_data.DirStringsOffset)
	}
	return vol_info, dirstrings
}

//Gets directory strings for all prefetch versions
func getDirectoryStrings(pf_file *os.File, volinfo_offset uint32, dirstrcnt uint32, dirstroff uint32) []string {
	dirstrings := []string{}
	cnt := 0
	count := dirstrcnt
	dirstr_offset := dirstroff
	pf_file.Seek(int64(volinfo_offset), 0)
	pf_file.Seek(int64(dirstr_offset), 1)

	for cnt < int(count) {
		str_len := make([]byte, 2)
		pf_file.Read(str_len)
		len := int(binary.LittleEndian.Uint16(str_len)) * 2
		dirstr := make([]byte, len)
		pf_file.Read(dirstr)
		data := strings.Replace(string(dirstr), "\x00", "", -1)
		null := make([]byte, 2)
		pf_file.Read(null)
		dirstrings = append(dirstrings, data)
		cnt += 1
	}
	return dirstrings
}

//Gets filename strings for all prefetch versions
func getFileNameStrings(pf_file *os.File, fnstrsize uint32, fnstroff uint32) []string {
	fnstrings := []string{}
	fnames := make([]byte, fnstrsize)
	pf_file.Seek(int64(fnstroff), 0)
	pf_file.Read(fnames)
	fns := formatFileNames(string(fnames))

	for _, f := range fns {
		if f != "" {
			fnstrings = append(fnstrings, strings.Replace(f, "\x00", "", -1))
		}
	}
	return fnstrings

}

//Opens an individual prefetch file
func openPrefetchFile(pf string) *os.File {
	pf_handle, err := os.Open(pf)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	return pf_handle
}

//Closes an individual prefetch file
func closePrefetchFile(pf string, open_file *os.File) {
	open_file.Close()
}

//Reads bytes from the specified start to end offset within a file
func readBytes(file *os.File, start int, end int) []byte {
	count := end - start
	bytes := make([]byte, count)

	_, err := file.ReadAt(bytes, int64(start))
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

//Gets the path for all prefetch files in a target path
func getPrefetchFiles(target_path string) []string {
	pf_files := []string{}

	err := filepath.Walk(target_path, func(path string, finfo os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !finfo.IsDir() {
			r, err := regexp.MatchString(".pf", finfo.Name())
			if err == nil && r {
				pf_files = append(pf_files, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return pf_files
}

//Formats an executable name
func formatExeName(exe_name []byte) string {
	tmp_name := strings.Split(string(exe_name), "\x00\x00")[0]
	fmat_exe_name := strings.Replace(string(tmp_name), "\x00", "", -1)
	return fmat_exe_name
}

//Formats the prefetch file version
func formatVersion(version []byte) int {
	fmat_version, _ := utf8.DecodeRuneInString(string(version))
	return int(fmat_version)
}

//Formats filename strings returned from the prefetch file
func formatFileNames(fnames string) []string {
	fmat_fnames := []string{}
	tmp_names := strings.Split(fnames, "\x00\x00")
	for _, fn := range tmp_names {
		if len(fn) > 0 {
			fmat_fname := strings.Replace(string(fn), "\x00", "", -1)
			fmat_fnames = append(fmat_fnames, fmat_fname)
		}
	}

	return fmat_fnames
}

//Format a Win32 timstamp into human readable
func formatWin32FileTime(ft uint64) time.Time {
	diff2unix := int64(11644473600000 * 10000)
	unixtime := int64(ft) - diff2unix
	return time.Unix(unixtime/10000000, 0).UTC()
}

//Potentially usable code at a later point or in a different project
// func (hdr *PrefetchHeader) MarshalJSON() []byte {
// 	type Alias PrefetchHeader
// 	j, err := json.MarshalIndent(&struct {
// 		Signature string
// 		Exename   string
// 		Version   int
// 		Filesize  uint32
// 		Hash      string
// 		*Alias
// 	}{
// 		Signature: string(hdr.Signature[:]),
// 		Exename:   formatExeName(hdr.Exename[:]),
// 		Version:   formatVersion(hdr.Version[:]),
// 		//Filesize:  formatFileSize(hdr.Filesize[:]),
// 		Filesize: binary.LittleEndian.Uint32(hdr.Filesize[:]),
// 		Hash:     hex.EncodeToString(hdr.Hash[:]),
// 		Alias:    (*Alias)(hdr),
// 	}, "", "    ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return j
// }

// func unMarshalJSON(data []byte) error {
// 	var headers []PrefetchHeader
// 	s, err := json.Unmarshal(data, &headers)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(s))
// 	return nil
// }
