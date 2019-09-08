package armeria

import "go.uber.org/zap"

// RegisterGameCommands registers all of the slash commands with the command manager.
func RegisterGameCommands() {
	commands := []*Command{
		{
			Name:    "commands",
			Help:    "Display top-level commands you have access to.",
			Handler: handleCommandsCommand,
		},
		{
			Name: "login",
			Help: "Log your character into the game world.",
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
			Name:   "logintoken",
			Help:   "Log your character into the game world (with a token).",
			Hidden: true,
			Permissions: &CommandPermissions{
				RequireNoCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:  "token",
					NoLog: true,
				},
			},
			Handler: handleLoginCommand,
		},
		{
			Name:   "create",
			Help:   "Create a new character.",
			Hidden: true,
			Permissions: &CommandPermissions{
				RequireNoCharacter: true,
			},
			Handler: handleCreateCommand,
		},
		{
			Name: "look",
			Help: "Look at something.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:             "at",
					Optional:         true,
					IncludeRemaining: true,
				},
			},
			Handler: handleLookCommand,
		},
		{
			Name: "glance",
			Help: "Glance around.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleGlanceCommand,
		},
		{
			Name: "say",
			Help: "Say something to everyone in your current room.",
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
			Help: "Move your character into a connecting room.",
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
			Help: "Manage rooms and their properties.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "edit",
					Help: "Open the editor panel for the current, or specified, room.",
					Arguments: []*CommandArgument{
						{
							Name:             "target",
							IncludeRemaining: true,
							Optional:         true,
						},
					},
					Handler: handleRoomEditCommand,
				},
				{
					Name: "set",
					Help: "Set attributes for the current, or specified, room. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "target",
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
					Handler: handleRoomSetCommand,
				},
				{
					Name: "create",
					Help: "Create a new room in the specified direction.",
					Arguments: []*CommandArgument{
						{
							Name: "direction",
						},
					},
					Handler: handleRoomCreateCommand,
				},
				{
					Name: "destroy",
					Help: "Destroy a room in the specified direction.",
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
			Help: "Manage characters.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_CHAREDIT",
			},
			Subcommands: []*Command{
				{
					Name: "list",
					Help: "List the characters in the game, optionally using a filter.",
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
					Help: "Set an attribute on the specified character. Leave value empty to revert to default.",
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
					Help: "Open the editor panel for the specified character.",
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
			Help: "Write the in-memory game data to disk.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_SYSOP",
			},
			Handler: handleSaveCommand,
		},
		{
			Name: "refresh",
			Help: "Re-render the data on the client.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleRefreshCommand,
		},
		{
			Name:     "whisper",
			AltNames: []string{"w"},
			Help:     "Send a private message to the specified character. They must be online.",
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
			Help:     "Reply to the last whisper you received.",
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
			Help: "Display a list of all characters who are currently online.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleWhoCommand,
		},
		{
			Name: "mob",
			Help: "Manage mobiles (npcs/monsters).",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "list",
					Help: "List the mobs in the game, optionally using a filter.",
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
					Help: "Create a new mob in the game.",
					Arguments: []*CommandArgument{
						{
							Name: "name",
						},
					},
					Handler: handleMobCreateCommand,
				},
				{
					Name: "edit",
					Help: "Edit a mob within the object editor.",
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
					Help: "Set a mob attribute. Leave value empty to revert to default.",
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
					Help: "Spawn a mob in your current room.",
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
				{
					Name: "iedit",
					Help: "Edit a specific mob instance within the object editor.",
					Arguments: []*CommandArgument{
						{
							Name: "uuid",
						},
					},
					Handler: handleMobInstanceEditCommand,
				},
				{
					Name: "iset",
					Help: "Set a mob instance attribute. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "uuid",
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
					Handler: handleMobInstanceSetCommand,
				},
			},
		},
		{
			Name: "area",
			Help: "Manage areas.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "create",
					Help: "Create a new area in the game.",
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
					Help: "List the areas in the game, optionally using a filter.",
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
					Help: "Open the editor panel for the specified area.",
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
			Help: "Manage items.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "list",
					Help: "List the items in the game, optionally using a filter.",
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
					Help: "Create a new item in the game.",
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
					Help: "Spawn an item in your current room.",
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
					Help: "Edit an item within the object editor.",
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
					Help: "Set an item attribute. Leave value empty to revert to default.",
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
					Name: "iedit",
					Help: "Edit a specific item instance within the object editor.",
					Arguments: []*CommandArgument{
						{
							Name: "uuid",
						},
					},
					Handler: handleItemInstanceEditCommand,
				},
				{
					Name: "iset",
					Help: "Set an item instance attribute. Leave value empty to revert to default.",
					Arguments: []*CommandArgument{
						{
							Name: "uuid",
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
					Handler: handleItemInstanceSetCommand,
				},
				{
					Name: "instances",
					Help: "View the locations of an item.",
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
			Help: "Wipe ALL objects in your current room.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Handler: handleWipeCommand,
		},
		{
			Name: "ghost",
			Help: "Bypass movement restrictions while moving.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_GHOST",
			},
			Handler: handleGhostCommand,
		},
		{
			Name: "password",
			Help: "Set a new password for your character.",
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
			Help:     "Teleport to the specified character or room.",
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
			Help:     "Copy and paste object attributes.",
			Permissions: &CommandPermissions{
				RequireCharacter:  true,
				RequirePermission: "CAN_BUILD",
			},
			Subcommands: []*Command{
				{
					Name: "copy",
					Help: "Copy attributes to your clipboard.",
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
					Help: "Paste attributes onto a matching object type.",
					Arguments: []*CommandArgument{
						{
							Name: "name",
						},
					},
					Handler: handleClipboardPasteCommand,
				},
				{
					Name:    "clear",
					Help:    "Clear the clipboard.",
					Handler: handleClipboardClearCommand,
				},
			},
		},
		{
			Name: "get",
			Help: "Grab an item from the ground.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:             "item",
					IncludeRemaining: true,
				},
			},
			Handler: handleGetCommand,
		},
		{
			Name: "drop",
			Help: "Drop an item onto the ground.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:             "item",
					IncludeRemaining: true,
				},
			},
			Handler: handleDropCommand,
		},
		{
			Name: "swap",
			Help: "Swap items in your inventory.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name: "source",
				},
				{
					Name: "destination",
				},
			},
			Handler: handleSwapCommand,
		},
		{
			Name: "autologin",
			Help: "Toggle auto-login for your character.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleAutoLoginCommand,
		},
		{
			Name: "channel",
			Help: "Join, leave or list talking channels you can participate in.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Subcommands: []*Command{
				{
					Name:    "list",
					Help:    "List channels you can join and participate in.",
					Handler: handleChannelListCommand,
				},
				{
					Name: "join",
					Help: "Join a channel.",
					Arguments: []*CommandArgument{
						{
							Name: "channel",
						},
					},
					Handler: handleChannelJoinCommand,
				},
				{
					Name: "say",
					Help: "Say someting to a channel.",
					Arguments: []*CommandArgument{
						{
							Name: "channel",
						},
						{
							Name:             "text",
							IncludeRemaining: true,
						},
					},
					Handler: handleChannelSayCommand,
				},
			},
		},
		{
			Name:     "settings",
			AltNames: []string{"setting"},
			Help:     "View or change settings for your character.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:     "name",
					Optional: true,
				},
				{
					Name:     "value",
					Optional: true,
				},
			},
			Handler: handleSettingsCommand,
		},
		{
			Name:     "bug",
			AltNames: []string{"idea"},
			Help:     "Report a bug or suggest an idea.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name:             "bug",
					Optional:         false,
					IncludeRemaining: true,
				},
			},
			Handler: handleBugCommand,
		},
		{
			Name: "give",
			Help: "Give an item to someone or something.",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Name: "target",
				},
				{
					Name:             "item",
					IncludeRemaining: true,
				},
			},
			Handler: handleGiveCommand,
		},
	}

	for _, cmd := range commands {
		Armeria.commandManager.RegisterCommand(cmd)
	}

	Armeria.log.Info("commands registered", zap.Int("count", len(commands)))
}
