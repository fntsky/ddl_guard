package ddlcmd

import "github.com/spf13/cobra"
var (
	runCmd =&cobra.Command{
		Use: "run",
		Short: "run server",
		Long: "run server",
		Run: func(_*cobra.Command,_[]string){
			runApp()
		},
	}
)

func Execute(){
	err:= runCmd.Execute()
	if err!=nil{
		panic(err)
	}
}