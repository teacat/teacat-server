package diskutil

// https://gist.github.com/lunny/9828326
import "syscall"

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func Usage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func (d DiskStatus) UsedGB() float64 {
	return float64(d.Used) / float64(GB)
}

func (d DiskStatus) FreeGB() float64 {
	return float64(d.Free) / float64(GB)
}

func (d DiskStatus) AllGB() float64 {
	return float64(d.All) / float64(GB)
}

func (d DiskStatus) UsedPercentage() float64 {
	return float64(float64(d.Used) / float64(d.All) * 100)
}
