package command

import (
	"errors"
	"fmt"
)

type Option struct {
	Name        string
	Description string
}

func (p *Option) Check() error {
	if p.Name == "" {
		return errors.New("option name required")
	}
	if p.Description == "" {
		return errors.New("option description required")
	}
	return nil
}

type DynamicParam struct {
	Func        func(input string) []string
	Flag        string
	Description string
}

func (dc *DynamicParam) Check() error {
	if dc.Flag == "" {
		return errors.New("dynamic param flag required")
	}
	if dc.Description == "" {
		return errors.New("dynamic param description required")
	}
	if dc.Func == nil {
		return errors.New("dynamic param fun required")
	}
	return nil
}

type Command struct {
	Name         string
	Description  string
	Commands     []*Command
	Options      []*Option
	DynamicParam *DynamicParam
	Run          func(cmd *ExecCmd)
}

func (cl *Command) AddCommand(cmds ...*Command) {
	for _, cmd := range cmds {
		cl.Commands = append(cl.Commands, cmd)
	}
}

func (cl *Command) Help() {
	clearLine() // 清除光标所在位置后的一行的标准输出
	fmt.Println(cl.Name + ":")
	if len(cl.Commands) > 0 {
		for _, subCmd := range cl.Commands {
			fmt.Printf("  %s\t%s\n", subCmd.Name, subCmd.Description)
		}
	} else {
		if cl.DynamicParam != nil {
			fmt.Printf("  (%s)\t%s\n", cl.DynamicParam.Flag, cl.DynamicParam.Description)
		}
		if len(cl.Options) > 0 {
			fmt.Println("  options:")
			for _, option := range cl.Options {
				fmt.Printf("    %s\t%s\n", option.Name, option.Description)
			}
		}
	}
}

func clearLine() {
	// 创建并打印包含光标移动到行首和清除行的控制字符序列
	fmt.Print("\r\033[K")
}

func (cl *Command) Check() error {
	if cl.Name == "" {
		return errors.New("command name required")
	}
	if cl.Description == "" {
		return errors.New("command description required")
	}
	if cl.Run == nil {
		if len(cl.Commands) == 0 {
			return errors.New("command func required")
		}
		// 如果还有子命令，则当前命令的Run默认为help命令
		cl.Run = func(cmd *ExecCmd) {
			cmd.Command.Help()
		}
	}
	if cl.DynamicParam != nil {
		if err := cl.DynamicParam.Check(); err != nil {
			return err
		}
	}
	if len(cl.Options) > 0 {
		for _, arg := range cl.Options {
			if err := arg.Check(); err != nil {
				return err
			}
		}
	}
	if len(cl.Commands) > 0 {
		for _, subCmd := range cl.Commands {
			if err := subCmd.Check(); err != nil {
				return err
			}
		}
	}
	return nil
}
