package pic

//字节转成png图片
import (
	"bytes"
	"image"
	"image/png"
)

func B2P(b []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// png图片转成base64字符串
func P2BASE64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

// 字节转成base64字符串
func B2BASE64(b []byte) (string, error) {
	img, err := B2P(b)
	if err != nil {
		return "", err
	}

	base64Str, err := P2BASE64(img)
	if err != nil {
		return "", err
	}
	return base64Str, nil
}
