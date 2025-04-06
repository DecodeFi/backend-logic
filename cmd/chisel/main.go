package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DecodeFi/backend-logic/internal/alchemy"
	"github.com/DecodeFi/backend-logic/internal/evm_inspect"
)

var cliName string = "chisel"

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func printUnknown(text string) {
	fmt.Println(text, ": command not found")
}

func displayHelp() {
	fmt.Printf(
		"Welcome to %v! These are the available commands: \n",
		cliName,
	)
	fmt.Println(".help    - Show available commands")
	fmt.Println(".clear   - Clear the terminal screen")
	fmt.Println(".exit    - Exits terminal ")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func handleInvalidCmd(text string) {
	defer printUnknown(text)
}

func cleanInput(text string) string {
	output := strings.TrimSpace(text)
	output = strings.ToLower(output)
	return output
}

type CommandHandler struct {
	alchemyClient    *alchemy.AlchemyClient
	evmInspectClient *evm_inspect.EvmInspectClient
}

func (c *CommandHandler) handleCmd(text string) {
	args := strings.Split(text, " ")

	cmd := args[0]

	if cmd == "block_number" {
		resp, _ := c.alchemyClient.BlockNumber(alchemy.NewBlockNumberRequest())
		fmt.Println(resp.Number)
	} else if cmd == "block" {
		block_no := args[1]
		resp, _ := c.alchemyClient.Block(alchemy.NewBlockRequest(block_no))
		fmt.Println(resp.Result.Hash)
		fmt.Println(resp.Result.Transactions[0])
	} else if cmd == "trace_block" {
		block_no := args[1]
		resp, _ := c.evmInspectClient.TraceBlock(block_no)
		fmt.Println(resp[0])
	} else {
		handleInvalidCmd(cmd)
	}
}

func main() {

	hanlder := &CommandHandler{
		alchemyClient:    alchemy.NewAlchemyClient(),
		evmInspectClient: evm_inspect.NewEvmInspectClient(),
	}

	// Hardcoded repl commands
	commands := map[string]interface{}{
		".help":  displayHelp,
		".clear": clearScreen,
	}

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		text := cleanInput(reader.Text())
		if command, exists := commands[text]; exists {
			command.(func())() // Call a hardcoded function
		} else if strings.EqualFold(".exit", text) {
			return
		} else {
			hanlder.handleCmd(text)
		}
		printPrompt()
	}

	fmt.Println()
}
