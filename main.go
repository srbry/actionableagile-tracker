package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	trackerAuth "github.com/srbry/actionableagile-tracker/tracker"

	tracker "github.com/xoebus/go-tracker"
)

type csvData [][]string

var csvHeader = []string{
	"Story ID",
	"Link",
	"Backlog",
	"Planned",
	"Started",
	"Finished",
	"Delivered",
	"Accepted",
}

type TransitionDates struct {
	Unscheduled string
	Planned     string
	Started     string
	Finished    string
	Delivered   string
	Accepted    string
	Rejected    string
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	fmt.Println("")

	fmt.Print("Enter Project Name: ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	config := trackerAuth.Config{
		Username: username,
		Password: password,
	}

	authClient := trackerAuth.NewClient(config)
	project, err := authClient.Project(projectName)
	if err != nil {
		log.Fatalf("Error getting project ID, %s", err.Error())
	}

	client := tracker.NewClient(authClient.Config.Token)
	projectClient := client.InProject(project.ID)

	stories, _, err := projectClient.Stories(tracker.StoriesQuery{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var outputData csvData
	outputData = append(outputData, csvHeader)
	for _, story := range stories {
		var outputRow []string
		outputRow = append(outputRow, strconv.Itoa(story.ID))
		outputRow = append(outputRow, story.URL)
		transitions, _ := authClient.StoryTransitions(project.ID, story.ID)
		var transitionDates TransitionDates
		for _, transition := range transitions {
			if transition.State == "unscheduled" {
				transitionDates.Unscheduled = transition.OccurredAt
			}
			if transition.State == "planned" {
				transitionDates.Planned = transition.OccurredAt
			}
			if transition.State == "started" {
				transitionDates.Started = transition.OccurredAt
			}
			if transition.State == "finished" {
				transitionDates.Finished = transition.OccurredAt
			}
			if transition.State == "delivered" {
				transitionDates.Delivered = transition.OccurredAt
			}
			if transition.State == "accepted" {
				transitionDates.Accepted = transition.OccurredAt
			}
		}
		outputRow = append(outputRow, transitionDates.Unscheduled)
		outputRow = append(outputRow, transitionDates.Planned)
		outputRow = append(outputRow, transitionDates.Started)
		outputRow = append(outputRow, transitionDates.Finished)
		outputRow = append(outputRow, transitionDates.Delivered)
		outputRow = append(outputRow, transitionDates.Accepted)
		outputData = append(outputData, outputRow)
	}

	var outputCSV = "actionableagile_import.csv"

	file, _ := os.Create(outputCSV)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range outputData {
		writer.Write(value)
	}
	fmt.Printf("Data written too %s\n", outputCSV)
}
