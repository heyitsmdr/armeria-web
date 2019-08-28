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
		p.client.ShowColorizedText("The picture could not be uploaded to a problem that occurred on the server.", ColorError)
		return
	}

	objectType := o["objectType"].(string)
	name := o["name"].(string)

	var oldKey string
	var editorData *ObjectEditorData
	switch objectType {
	case "character":
		c := Armeria.characterManager.CharacterByName(name)
		oldKey = c.Attribute(AttributePicture)
		_ = c.SetAttribute(AttributePicture, k)
		editorData = c.EditorData()
		p.client.ShowColorizedText(
			fmt.Sprintf("A picture has been uploaded and set for character %s.", c.FormattedName()),
			ColorSuccess,
		)
		if c.Online() {
			for _, chars := range c.Room().Here().Characters(true, nil) {
				chars.Player().client.SyncRoomObjects()
			}
		}
	case "mob":
		m := Armeria.mobManager.MobByName(name)
		oldKey = m.Attribute(AttributePicture)
		m.SetAttribute(AttributePicture, k)
		editorData = m.EditorData()
		p.client.ShowColorizedText(
			fmt.Sprintf("A picture has been uploaded and set for mob %s.", TextStyle(m.Name(), TextStyleBold)),
			ColorSuccess,
		)
		for _, mi := range m.Instances() {
			for _, chars := range mi.Room().Here().Characters(true, nil) {
				chars.Player().client.SyncRoomObjects()
			}
		}
	case "item":
		i := Armeria.itemManager.ItemByName(name)
		oldKey = i.Attribute(AttributePicture)
		i.SetAttribute(AttributePicture, k)
		editorData = i.EditorData()
		p.client.ShowColorizedText(
			fmt.Sprintf("A picture has been uploaded and set for item %s.", TextStyle(i.Name(), TextStyleBold)),
			ColorSuccess,
		)
		for _, ii := range i.Instances() {
			ctr := Armeria.registry.GetObjectContainer(ii.ID())
			if ctr.ParentType() == ContainerParentTypeRoom {
				for _, chars := range ii.Room().Here().Characters(true, nil) {
					chars.Player().client.SyncRoomObjects()
				}
			} else if ctr.ParentType() == ContainerParentTypeCharacter {
				if ii.Character().Online() {
					ii.Character().Player().client.SyncInventory()
				}
			}
		}
	default:
		p.client.ShowColorizedText("The picture was uploaded as an invalid type.", ColorError)
		DeleteObjectPictureFromDisk(k)
		return
	}

	if oldKey != k && len(oldKey) > 0 {
		DeleteObjectPictureFromDisk(oldKey)
	}

	editorOpen := p.Character().TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		p.client.ShowObjectEditor(editorData)
	}
}

// SaveObjectPictureToDisk stores an object picture on the disk and returns the key.
func SaveObjectPictureToDisk(o map[string]interface{}) string {
	objectType := o["objectType"].(string)
	name := o["name"].(string)
	pictureType := o["pictureType"].(string)
	pictureData := o["pictureData"].(string)

	hash := md5.Sum([]byte(pictureData))
	normalizedName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	key := fmt.Sprintf("%s-%s-%x", objectType, normalizedName, hash)

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

	_ = f.Sync()

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
