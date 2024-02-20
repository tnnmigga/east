package msgbus

type castOpt struct {
	key   string
	value any
}

const (
	keyNone         = "none"
	keyNonuseStream = "nonuse-stream"
	keyOneOfMods    = "one-of-mods"
	keyServerID     = "server-id"
)

func findCastOpt[T any](opts []castOpt, key string) (value T, find bool) {
	for _, opt := range opts {
		if opt.key == key {
			return opt.value.(T), true
		}
	}
	return value, false
}

func NonuseStream() castOpt {
	return castOpt{
		key:   keyNonuseStream,
		value: true,
	}
}

func OneOfMods(modName string) castOpt {
	return castOpt{
		key:   keyOneOfMods,
		value: modName,
	}
}

func ServerID(serverID uint32) castOpt {
	return castOpt{
		key:   keyServerID,
		value: serverID,
	}
}
