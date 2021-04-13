package main

import (
	"fmt"
	"io"
)

type jattach struct {
	LocalPath            string `env:"JATTACH_LOCAL_PATH" envDefault:"./jattach"`
	RemotePath           string `env:"JATTACH_REMOTE_PATH" envDefault:"/tmp/jattach"`
	OutputDumpPath       string `env:"JATTACH_OUT_DUMP_PATH" envDefault:"/tmp/jattach-heapdump.hprof"`
	DumpScriptRemotePath string `env:"JATTACH_DUMP_SCRIPT_PATH" envDefault:"/tmp/jattach-dump.sh"`
}

const (
	dumpScriptLocalPath = "/tmp/jattach-dump.sh"
	pulledDumpLocalPath = "/tmp/jattach-heapdump.hprof"
	localJAttachPath    = "/jattach"
)

func (j jattach) WriteDumpScript(writer io.Writer) error {
	commands := []string{
		fmt.Sprintf("chmod +x %s", j.RemotePath),
		fmt.Sprintf("%s $(pidof -s java) dumpheap %s", j.RemotePath, j.OutputDumpPath),
	}

	for _, c := range commands {
		_, err := writer.Write([]byte(c + "\n"))
		if err != nil {
			return fmt.Errorf("write command: '%s': %w", c, err)
		}
	}
	return nil
}
