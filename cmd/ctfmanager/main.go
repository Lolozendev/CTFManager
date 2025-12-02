package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Lolozendev/CTFManager/internal/app/challenge"
	"github.com/Lolozendev/CTFManager/internal/app/compose"
	"github.com/Lolozendev/CTFManager/internal/app/team"
	"github.com/Lolozendev/CTFManager/internal/config"
	"github.com/Lolozendev/CTFManager/internal/logger"
	"github.com/Lolozendev/CTFManager/internal/model"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
	log = logger.Get()
)

func main() {
	defer logger.Close()

	// Initialize configuration
	cfg = config.Default()

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "ctfmanager",
		Short: "CTFManager - Manage dockerized CTF environments",
		Long: `CTFManager is a CLI tool for managing Docker-based CTF (Capture The Flag) environments.
It helps you create and manage teams, challenges, and their associated infrastructure.`,
	}

	// Add subcommands
	rootCmd.AddCommand(setupCmd())
	rootCmd.AddCommand(teamCmd())
	rootCmd.AddCommand(challengeCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

// setupCmd initializes the CTF environment
func setupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Initialize the CTF environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Setting up CTF environment...")

			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
			}

			log.Info("CTF environment setup complete!")
			return nil
		},
	}
}

// teamCmd returns the team management command
func teamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "Manage CTF teams",
	}

	cmd.AddCommand(teamListCmd())
	cmd.AddCommand(teamCreateCmd())
	cmd.AddCommand(teamDeleteCmd())
	cmd.AddCommand(teamEnableCmd())
	cmd.AddCommand(teamDisableCmd())

	return cmd
}

func teamListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all teams",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := team.New(cfg, log)
			teams, err := mgr.List()
			if err != nil {
				return err
			}

			if len(teams) == 0 {
				log.Info("No teams found")
				return nil
			}

			fmt.Println("\nTeams:")
			for _, t := range teams {
				status := "enabled"
				if !t.Enabled {
					status = "disabled"
				}
				fmt.Printf("  [%d] %s (%s)\n", t.ID, t.Name, status)
			}
			fmt.Println()

			return nil
		},
	}
}

func teamCreateCmd() *cobra.Command {
	var members []string

	cmd := &cobra.Command{
		Use:   "create <id> <name>",
		Short: "Create a new team",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid team ID: %w", err)
			}

			name := args[1]

			mgr := team.New(cfg, log)
			if err := mgr.Create(id, name, members); err != nil {
				return err
			}

			// Generate compose file
			challMgr := challenge.New(cfg, log)
			challenges, err := challMgr.ListEnabled()
			if err != nil {
				return fmt.Errorf("failed to list challenges: %w", err)
			}

			composeGen := compose.New(cfg, log)

			teamModel := model.Team{
				ID:      id,
				Name:    name,
				Members: stringSliceToMembers(members),
				Enabled: true,
			}

			composeYAML, err := composeGen.Generate(teamModel, challenges)
			if err != nil {
				return fmt.Errorf("failed to generate compose file: %w", err)
			}

			// Write compose file
			composePath := cfg.GetTeamPath(model.FormatChallengeName(id, name, true)) + "/compose.yml"
			if err := os.WriteFile(composePath, []byte(composeYAML), 0644); err != nil {
				return fmt.Errorf("failed to write compose file: %w", err)
			}

			fmt.Printf("\n✓ Team '%s' created successfully (ID: %d)\n", name, id)
			fmt.Printf("  Compose file: %s\n\n", composePath)

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&members, "members", "m", []string{}, "Team members (comma-separated)")

	return cmd
}

func teamDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := team.New(cfg, log)
			if err := mgr.Delete(args[0]); err != nil {
				return err
			}

			fmt.Printf("\n✓ Team '%s' deleted successfully\n\n", args[0])
			return nil
		},
	}
}

func teamEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <name> <id>",
		Short: "Enable a disabled team",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid team ID: %w", err)
			}

			mgr := team.New(cfg, log)
			if err := mgr.Enable(args[0], id); err != nil {
				return err
			}

			fmt.Printf("\n✓ Team '%s' enabled successfully\n\n", args[0])
			return nil
		},
	}
}

func teamDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <name>",
		Short: "Disable a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := team.New(cfg, log)
			if err := mgr.Disable(args[0]); err != nil {
				return err
			}

			fmt.Printf("\n✓ Team '%s' disabled successfully\n\n", args[0])
			return nil
		},
	}
}

// challengeCmd returns the challenge management command
func challengeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "challenge",
		Short: "Manage CTF challenges",
	}

	cmd.AddCommand(challengeListCmd())
	cmd.AddCommand(challengeValidateCmd())
	cmd.AddCommand(challengeEnableCmd())
	cmd.AddCommand(challengeDisableCmd())

	return cmd
}

func challengeListCmd() *cobra.Command {
	var showDisabled bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all challenges",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := challenge.New(cfg, log)
			challenges, err := mgr.List()
			if err != nil {
				return err
			}

			if len(challenges) == 0 {
				log.Info("No challenges found")
				return nil
			}

			fmt.Println("\nChallenges:")
			for _, ch := range challenges {
				if !showDisabled && !ch.Enabled {
					continue
				}

				status := "enabled"
				if !ch.Enabled {
					status = "disabled"
				}

				networkID := "N/A"
				if ch.Enabled {
					networkID = fmt.Sprintf("%d", ch.NetworkID)
				}

				fmt.Printf("  [%s] %s (%s)\n", networkID, ch.Name, status)
			}
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&showDisabled, "all", "a", false, "Show disabled challenges")

	return cmd
}

func challengeValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate all challenges",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := challenge.New(cfg, log)
			if err := mgr.Validate(); err != nil {
				return err
			}

			fmt.Println("\n✓ All challenges are valid\n")
			return nil
		},
	}
}

func challengeEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <name> <network-id>",
		Short: "Enable a disabled challenge",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			networkID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid network ID: %w", err)
			}

			mgr := challenge.New(cfg, log)
			if err := mgr.Enable(args[0], networkID); err != nil {
				return err
			}

			fmt.Printf("\n✓ Challenge '%s' enabled with network ID %d\n\n", args[0], networkID)
			return nil
		},
	}
}

func challengeDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <name>",
		Short: "Disable a challenge",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr := challenge.New(cfg, log)
			if err := mgr.Disable(args[0]); err != nil {
				return err
			}

			fmt.Printf("\n✓ Challenge '%s' disabled successfully\n\n", args[0])
			return nil
		},
	}
}

// Helper functions
func stringSliceToMembers(names []string) []model.Member {
	members := make([]model.Member, len(names))
	for i, name := range names {
		members[i] = model.Member{Username: name}
	}
	return members
}
