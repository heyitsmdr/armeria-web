package armeria

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

// StoreObjectPicture handles the client-initiated process of storing an object picture.
func StoreObjectPicture(p *Player, o map[string]interface{}) {
	// TODO: check permissions

	k := SaveObjectPictureToDisk(o)
	if len(k) == 0 {
		p.clientActions.ShowColorizedText("The picture could not be uploaded to a problem that occurred on the server.", ColorError)
		return
	}

	objectType := o["objectType"].(string)
	name := o["name"].(string)

	var oldKey string
	var editorData *ObjectEditorData
	switch objectType {
	case "character":
		c := Armeria.characterManager.GetCharacterByName(name)
		oldKey = c.GetAttribute("picture")
		c.SetAttribute("picture", k)
		editorData = c.GetEditorData()
		p.clientActions.ShowColorizedText(
			fmt.Sprintf("A picture has been uploaded and set for %s.", c.GetFName()),
			ColorSuccess,
		)
		for _, chars := range p.GetCharacter().GetRoom().GetCharacters(nil) {
			chars.GetPlayer().clientActions.SyncRoomObjects()
		}
	default:
		p.clientActions.ShowColorizedText("The picture was uploaded as an invalid type.", ColorError)
		DeleteObjectPictureFromDisk(k)
		return
	}

	if oldKey != k && len(oldKey) > 0 {
		DeleteObjectPictureFromDisk(oldKey)
	}

	editorOpen := p.GetCharacter().GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		p.clientActions.ShowObjectEditor(editorData)
	}
}

// SaveObjectPictureToDisk stores an object picture on the disk and returns the key.
func SaveObjectPictureToDisk(o map[string]interface{}) string {
	objectType := o["objectType"].(string)
	name := o["name"].(string)
	pictureType := o["pictureType"].(string)
	pictureData := o["pictureData"].(string)

	hash := md5.Sum([]byte(pictureData))
	key := fmt.Sprintf("%s-%s-%x", objectType, strings.ToLower(name), hash)

	dec, err := base64.StdEncoding.DecodeString(pictureData)
	if err != nil {
		Armeria.log.Error("error decoding base64 picture upload",
			zap.String("type", objectType),
			zap.String("name", name),
			zap.Error(err),
		)
		return ""
	}

	var ext string
	switch pictureType {
	case "image/png":
		ext = "png"
	case "image/jpeg":
		ext = "jpg"
	case "image/jpg":
		ext = "jpg"
	default:
		ext = "png"
	}

	key = fmt.Sprintf("%s.%s", key, ext)

	pictureFile := fmt.Sprintf("%s/%s", Armeria.objectImagesPath, key)
	f, err := os.Create(pictureFile)
	if err != nil {
		Armeria.log.Error("error creating object picture file on disk",
			zap.String("file", pictureFile),
			zap.Error(err),
		)
		return ""
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		Armeria.log.Error("error writing object picture file to disk",
			zap.String("file", pictureFile),
			zap.Error(err),
		)
		return ""
	}
	f.Sync()

	Armeria.log.Info("wrote object picture to disk",
		zap.String("file", pictureFile),
	)

	return key
}

// DeleteObjectPictureFromDisk removes an object picture from the disk.
func DeleteObjectPictureFromDisk(k string) {
	pictureFile := fmt.Sprintf("%s/%s", Armeria.objectImagesPath, k)
	err := os.Remove(pictureFile)
	if err != nil {
		Armeria.log.Error("error removing old object picture",
			zap.String("file", pictureFile),
			zap.Error(err),
		)
		return
	}

	Armeria.log.Info("removed old object picture from disk",
		zap.String("file", pictureFile),
	)
}
