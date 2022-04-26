package cmd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	. "github.com/bhojpur/deploy/pkg/entities"

	tablewriter "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type EntityDifference struct {
	OriginalEntity Entity `json:"originalEntity,omitempty" yaml:"originalEntity,omitempty"`
	TargetEntity   Entity `json:"targetEntity,omitempty" yaml:"targetEntity,omitempty"`
	Kind           string `json:"kind" yaml:"kind"`
	Descr          string `json:"descr,omitempty" yaml:"descr,omitempty"`
	Missing        bool   `json:"missing" yaml:"missing"`
}

func getCurrentStatus(store *EntitiesStore, usersFile, groupsFile, shadowFile, gshadowFile string) error {

	mUsers, err := ParseUser(usersFile)
	if err != nil {
		return err
	}

	mGroups, err := ParseGroup(groupsFile)
	if err != nil {
		return err
	}

	mShadows, err := ParseShadow(shadowFile)
	if err != nil {
		return err
	}

	mGShadows, err := ParseGShadow(gshadowFile)
	if err != nil {
		return err
	}

	store.Users = mUsers
	store.Groups = mGroups
	store.Shadows = mShadows
	store.GShadows = mGShadows

	return nil
}

func compare(currentStore, store *EntitiesStore, jsonOutput bool) error {

	differences := []EntityDifference{}

	// Check users: I check that all entities defined in the specs are available and equal.
	// Not in reverse.
	for name, u := range store.Users {
		cUser, ok := currentStore.GetUser(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: u,
				Missing:      true,
				Kind:         u.GetKind(),
				Descr:        fmt.Sprintf("User %s is not present.", name),
			})
			continue
		}

		if (u.Uid >= 0 && cUser.Uid != u.Uid) ||
			(u.Group == "" && cUser.Gid != u.Gid) ||
			cUser.Homedir != u.Homedir || cUser.Shell != u.Shell {
			differences = append(differences, EntityDifference{
				OriginalEntity: cUser,
				TargetEntity:   u,
				Missing:        false,
				Kind:           u.GetKind(),
				Descr:          fmt.Sprintf("User %s has difference.", name),
			})
		}
	}

	// Check groups
	for name, g := range store.Groups {
		cGroup, ok := currentStore.GetGroup(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: g,
				Missing:      true,
				Kind:         g.GetKind(),
				Descr:        fmt.Sprintf("Group %s is not present.", name),
			})
			continue
		}

		if cGroup.Password != g.Password ||
			(g.Gid != nil && *g.Gid >= 0 && cGroup.Gid != g.Gid) ||
			cGroup.Users != g.Users {
			differences = append(differences, EntityDifference{
				OriginalEntity: cGroup,
				TargetEntity:   g,
				Missing:        false,
				Kind:           g.GetKind(),
				Descr:          fmt.Sprintf("Group %s has difference.", name),
			})
		}
	}

	// Check shadow
	for name, s := range store.Shadows {
		cShadow, ok := currentStore.GetShadow(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: s,
				Missing:      true,
				Kind:         s.GetKind(),
				Descr:        fmt.Sprintf("Shadow with username %s is not present.", name),
			})
			continue
		}

		if cShadow.MinimumChanged != s.MinimumChanged ||
			cShadow.MaximumChanged != s.MaximumChanged ||
			cShadow.Warn != s.Warn ||
			cShadow.Inactive != s.Inactive ||
			cShadow.Expire != s.Expire {
			differences = append(differences, EntityDifference{
				OriginalEntity: cShadow,
				TargetEntity:   s,
				Missing:        false,
				Kind:           s.GetKind(),
				Descr:          fmt.Sprintf("Shadow with user %s has difference.", name),
			})
		}
	}

	// Check gshadow
	for name, s := range store.GShadows {
		cGShadow, ok := currentStore.GetGShadow(name)
		if !ok {
			differences = append(differences, EntityDifference{
				TargetEntity: s,
				Missing:      true,
				Kind:         s.GetKind(),
				Descr:        fmt.Sprintf("GShadow with name %s is not present.", name),
			})
			continue
		}

		if cGShadow.Password != s.Password ||
			cGShadow.Administrators != s.Administrators ||
			cGShadow.Members != s.Members {
			differences = append(differences, EntityDifference{
				OriginalEntity: cGShadow,
				TargetEntity:   s,
				Missing:        false,
				Kind:           s.GetKind(),
				Descr:          fmt.Sprintf("GShadow with name %s has difference.", name),
			})
		}
	}

	if jsonOutput {
		data, _ := json.Marshal(differences)
		fmt.Println(string(data))
	} else {

		table := tablewriter.NewWriter(os.Stdout)
		table.SetBorders(tablewriter.Border{
			Left:   true,
			Top:    true,
			Right:  true,
			Bottom: true,
		})
		table.SetColWidth(50)
		table.SetHeader([]string{
			"Kind", "Name", "Missing", "Difference",
		})
		for _, d := range differences {
			var name string
			switch d.Kind {
			case UserKind:
				name = (d.TargetEntity.(UserPasswd)).Username
			case ShadowKind:
				name = (d.TargetEntity.(Shadow)).Username
			case GroupKind:
				name = (d.TargetEntity.(Group)).Name
			case GShadowKind:
				name = (d.TargetEntity.(GShadow)).Name
			}

			table.Append([]string{
				d.Kind,
				name,
				fmt.Sprintf("%v", d.Missing),
				d.Descr,
			})

		}

		table.Render()

	}

	return nil
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare entities present with specs.",
	Long: `
Compare entities of the system with the specs available in the specified directory.

To read /etc/shadow and /etc/gshadow requires root permissions.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")
		if len(specsdirs) == 0 {
			return errors.New("At least one specs directory is needed.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		specsdirs, _ := cmd.Flags().GetStringArray("specs-dir")
		usersFile, _ := cmd.Flags().GetString("users-file")
		groupsFile, _ := cmd.Flags().GetString("groups-file")
		shadowFile, _ := cmd.Flags().GetString("shadow-file")
		gShadowFile, _ := cmd.Flags().GetString("gshadow-file")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		store := NewEntitiesStore()
		currentStore := NewEntitiesStore()

		// Load sepcs
		for _, d := range specsdirs {
			err := store.Load(d)
			if err != nil {
				return errors.New(
					"Error on load specs from directory " + d + ": " + err.Error())
			}
		}

		// Retrieve current information
		err := getCurrentStatus(currentStore,
			usersFile, groupsFile, shadowFile, gShadowFile,
		)
		if err != nil {
			return errors.New(
				"Error on retrieve current entities status: " + err.Error(),
			)
		}

		err = compare(currentStore, store, jsonOutput)
		if err != nil {
			return errors.New(
				"Error on compare entities stores: " + err.Error(),
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)

	var flags = compareCmd.Flags()
	flags.StringArrayP("specs-dir", "s", []string{},
		"Define the directory where read entities specs. At least one directory is needed.")
	flags.String("users-file", UserDefault(""), "Define custom users file.")
	flags.String("groups-file", GroupsDefault(""), "Define custom groups file.")
	flags.String("shadow-file", ShadowDefault(""), "Define custom shadow file.")
	flags.String("gshadow-file", GShadowDefault(""), "Define custom gshadow file.")
	flags.Bool("json", false, "Show in JSON format.")
}
