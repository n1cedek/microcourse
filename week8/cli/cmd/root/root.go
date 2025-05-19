package root

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "my-app",
	Short: "Моё cli приложение",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Что-то создает!",
}
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Что-то удаялет!",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Создает нового пользователя!",
	Run: func(cmd *cobra.Command, args []string) {
		userNameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("Failed to get username: %s\n", err.Error())
		}
		log.Printf("User %s creared\n", userNameStr)
	},
}
var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Удаляет пользователя!",
	Run: func(cmd *cobra.Command, args []string) {
		userNameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("Failed to get username: %s\n", err.Error())
		}
		log.Printf("User %s deletes\n", userNameStr)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)

	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())

	}
}
