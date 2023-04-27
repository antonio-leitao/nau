package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	utils "github.com/antonio-leitao/nau/lib/utils"
	"github.com/sahilm/fuzzy"
)

type Projects []utils.Project

func (p Projects) String(i int) string {
	return p[i].Name
}
func (p Projects) Len() int {
	return len(p)
}

func compressAndMove(srcDir string, destDir string) error {
	// get the name of the source directory
	srcDirName := filepath.Base(srcDir)

	// tar + gzip
	var buf bytes.Buffer
	if err := compress(srcDir, &buf); err != nil {
		return err
	}

	// write the .tar.gzip to file
	compressedFileName := fmt.Sprintf("%s.tar.gzip", srcDirName)
	compressedFilePath := filepath.Join(destDir, compressedFileName)
	fileToWrite, err := os.OpenFile(compressedFilePath, os.O_CREATE|os.O_RDWR, os.FileMode(600))
	if err != nil {
		return err
	}
	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		return err
	}

	// delete the source directory
	err = os.RemoveAll(srcDir)
	if err != nil {
		return err
	}

	return nil
}

func compress(src string, buf io.Writer) error {
	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}

func runMakeArchive(targetDir string) error {
	// Change to the target directory
	err := os.Chdir(targetDir)
	if err != nil {
		return err
	}
	// Run the "make archive" command
	cmd := exec.Command("make", "archive")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("make archive failed: %s", output)
	}

	return nil
}

func archiveProject(destDir, srcDir string) {
	//before archiving run the make archive target on that directory.
	_ = runMakeArchive(srcDir)
	//archive the project
	err := compressAndMove(srcDir, destDir)
	if err != nil {
		panic(err)
	}
}

func Archive(config utils.Config, query string) {
	projectList, _ := utils.GetProjects(config)
	projects := Projects(projectList)
	candidates := fuzzy.FindFrom(query, projects)

	//exit it nothing is found
	if len(candidates) == 0 {
		fmt.Println("ERROR: No project found")
		os.Exit(1)
	}

	projectPath := projects[candidates[0].Index].Path
    archivesPath,err := utils.ConvertPath(config.Archives_path)
    if err !=nil{
        fmt.Println("Archives not found")
    }
	archiveProject(archivesPath,projectPath)
}
