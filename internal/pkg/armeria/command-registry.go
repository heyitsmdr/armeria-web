package armeria

import "go.uber.org/zap"

func RegisterGameCommands() {
	commands := []*Command{
		{
			Name:    "commands",
			Help:    "Displays top-level commands you have access to.",
			Handler: handleCommandsCommand,
		},
		{
			Name: "login",
			Help: "Logs your character into the game world.",
			Permissions: &CommandPermissions{
				RequireNoCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name: "character",
				},
				{
					Name:  "password",
					NoLog: true,
				},
			},
			Handler: handleLoginCommand,
		},
		{
			Name: "look",
			Help: "Looks at something.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleLookCommand,
		},
		{
			Name: "say",
			Help: "Says something to everyone in the same room as you.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:             "text",
					IncludeRemaining: true,
				},
			},
			Handler: handleSayCommand,
		},
		{
			Name: "move",
			Help: "Moves your character into a connecting room.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name: "direction",
				},
			},
			Handler: handleMoveCommand,
		},
		{Name: "north", Alias: "move north"},
		{Name: "south", Alias: "move south"},
		{Name: "east", Alias: "move east"},
		{Name: "west", Alias: "move west"},
		{Name: "up", Alias: "move up"},
		{Name: "down", Alias: "move down"},
		{
			Name: "room",
			Help: "Allows you to manage rooms.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name:    "edit",
					Help:    "Opens the editor panel for the current room.",
					Handler: handleRoomEditCommand,
				},
				{
					Name: "set",
					Help: "Allows you to set a room attribute. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "property",
						},
						{
							Name:             "value",
							IncludeRemaining: true,
							Optional:         true,
						},
					},
					Handler: handleRoomSetCommand,
				},
				{
					Name: "create",
					Help: "Creates a new room in a particular direction.",
					Arguments: []*CommandArgument{
						{
							Name: "direction",
						},
					},
					Handler: handleRoomCreateCommand,
				},
				{
					Name: "destroy",
					Help: "Destroys a room in a particular direction.",
					Arguments: []*CommandArgument{
						{
							Name: "direction",
						},
					},
					Handler: handleRoomDestroyCommand,
				},
			},
		},
		{
			Name: "character",
			Help: "Allows you to manage characters.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_CHAREDIT",
			},
			Subcommands: []*Command{
				{
					Name: "list",
					Help: "Lists the characters in the game, optionally using a filter.",
					Arguments: []*CommandArgument{
						{
							Name:     "filter",
							Optional: true,
						},
					},
					Handler: handleCharacterListCommand,
				},
				{
					Name: "set",
					Help: "Allows you to set a character attribute. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "character",
						},
						{
							Name: "property",
						},
						{
							Name:             "value",
							IncludeRemaining: true,
							Optional:         true,
						},
					},
					Handler: handleCharacterSetCommand,
				},
				{
					Name: "edit",
					Help: "Opens the editor panel for a character.",
					Arguments: []*CommandArgument{
						{
							Name:     "character",
							Optional: true,
						},
					},
					Handler: handleCharacterEditCommand,
				},
			},
		},
		{
			Name: "save",
			Help: "Writes the in-memory game data to disk.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_SYSOP",
			},
			Handler: handleSaveCommand,
		},
		{
			Name: "refresh",
			Help: "Asks the server to re-render the data on the client.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleRefreshCommand,
		},
		{
			Name:     "whisper",
			AltNames: []string{"w"},
			Help:     "Sends a private message to another online character.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name: "target",
				},
				{
					Name:             "message",
					IncludeRemaining: true,
				},
			},
			Handler: handleWhisperCommand,
		},
		{
			Name:     "reply",
			AltNames: []string{"r"},
			Help:     "Replies to the last whisper received.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:             "message",
					IncludeRemaining: true,
				},
			},
			Handler: handleReplyCommand,
		},
		{
			Name: "who",
			Help: "Displays the characters currently playing.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleWhoCommand,
		},
		{
			Name: "mob",
			Help: "Allows you to manage mobiles (npcs/monsters).",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "list",
					Help: "Lists the mobs in the game, optionally using a filter.",
					Arguments: []*CommandArgument{
						{
							Name:     "filter",
							Optional: true,
						},
					},
					Handler: handleMobListCommand,
				},
				{
					Name: "create",
					Help: "Creates a new mob in the game.",
					Arguments: []*CommandArgument{
						{
							Name: "name",
						},
					},
					Handler: handleMobCreateCommand,
				},
				{
					Name: "edit",
					Help: "Opens the editor panel for a mob.",
					Arguments: []*CommandArgument{
						{
							Name:             "mob",
							IncludeRemaining: true,
						},
					},
					Handler: handleMobEditCommand,
				},
				{
					Name: "set",
					Help: "Allows you to set a mob attribute. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "mob",
						},
						{
							Name: "property",
						},
						{
							Name:             "value",
							IncludeRemaining: true,
							Optional:         true,
						},
					},
					Handler: handleMobSetCommand,
				},
				{
					Name: "spawn",
					Help: "Spawns a mob at this location.",
					Arguments: []*CommandArgument{
						{
							Name:             "mob",
							IncludeRemaining: true,
						},
					},
					Handler: handleMobSpawnCommand,
				},
				{
					Name: "instances",
					Help: "View the locations of a particular mob.",
					Arguments: []*CommandArgument{
						{
							Name:             "mob",
							IncludeRemaining: true,
						},
					},
					Handler: handleMobInstancesCommand,
				},
			},
		},
		{
			Name: "area",
			Help: "Allows you to manage areas.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "create",
					Help: "Creates a new area in the game.",
					Arguments: []*CommandArgument{
						{
							Name:             "name",
							IncludeRemaining: true,
						},
					},
					Handler: handleAreaCreateCommand,
				},
				{
					Name: "list",
					Help: "Lists the areas in the game, optionally using a filter.",
					Arguments: []*CommandArgument{
						{
							Name:     "filter",
							Optional: true,
						},
					},
					Handler: handleAreaListCommand,
				},
				{
					Name: "edit",
					Help: "Opens the editor panel for an area.",
					Arguments: []*CommandArgument{
						{
							Name:             "area",
							Optional:         true,
							IncludeRemaining: true,
						},
					},
					Handler: handleAreaEditCommand,
				},
			},
		},
		{
			Name: "item",
			Help: "Allows you to manage items.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "list",
					Help: "Lists the items in the game, optionally using a filter.",
					Arguments: []*CommandArgument{
						{
							Name:     "filter",
							Optional: true,
						},
					},
					Handler: handleItemListCommand,
				},
				{
					Name: "create",
					Help: "Creates a new item in the game.",
					Arguments: []*CommandArgument{
						{
							Name:             "name",
							IncludeRemaining: true,
						},
					},
					Handler: handleItemCreateCommand,
				},
				{
					Name: "spawn",
					Help: "Spawns an item at this location.",
					Arguments: []*CommandArgument{
						{
							Name:             "item",
							IncludeRemaining: true,
						},
					},
					Handler: handleItemSpawnCommand,
				},
				{
					Name: "edit",
					Help: "Opens the editor panel for an item.",
					Arguments: []*CommandArgument{
						{
							Name:             "item",
							IncludeRemaining: true,
						},
					},
					Handler: handleItemEditCommand,
				},
				{
					Name: "set",
					Help: "Allows you to set an item attribute. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "item",
						},
						{
							Name: "property",
						},
						{
							Name:             "value",
							IncludeRemaining: true,
							Optional:         true,
						},
					},
					Handler: handleItemSetCommand,
				},
				{
					Name: "instances",
					Help: "View the locations of a particular item.",
					Arguments: []*CommandArgument{
						{
							Name:             "item",
							IncludeRemaining: true,
						},
					},
					Handler: handleItemInstancesCommand,
				},
			},
		},
		{
			Name: "wipe",
			Help: "Wipes the objects in the same room.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Handler: handleWipeCommand,
		},
		{
			Name: "ghost",
			Help: "Allows you to bypass restrictions when moving.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_GHOST",
			},
			Handler: handleGhostCommand,
		},
		{
			Name: "password",
			Help: "Sets a new password for your character.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:  "password",
					NoLog: true,
				},
			},
			Handler: handlePasswordCommand,
		},
		{
			Name:     "teleport",
			AltNames: []string{"tp"},
			Help:     "Teleports to a character or location.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_TELEPORT",
			},
			Arguments: []*CommandArgument{
				{
					Name:             "destination",
					IncludeRemaining: true,
				},
			},
			Handler: handleTeleportCommand,
		},
		{
			Name:     "clipboard",
			AltNames: []string{"cb"},
			Help:     "Allows you to copy and paste object attributes easily.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "copy",
					Help: "Copies one or more object attributes onto your clipboard.",
					Arguments: []*CommandArgument{
						{
							Name: "type",
						},
						{
							Name: "name",
						},
						{
							Name:             "attributes",
							IncludeRemaining: true,
						},
					},
					Handler: handleClipboardCopyCommand,
				},
				{
					Name: "paste",
					Help: "Pastes attributes onto a matching object type.",
					Arguments: []*CommandArgument{
						{
							Name: "name",
						},
					},
					Handler: handleClipboardPasteCommand,
				},
				{
					Name:    "clear",
					Help:    "Clears the clipboard.",
					Handler: handleClipboardClearCommand,
				},
			},
		},
	}

	for _, cmd := range commands {
		Armeria.commandManager.RegisterCommand(cmd)
	}

	Armeria.log.Info("commands registered", zap.Int("count", len(commands)))
}
