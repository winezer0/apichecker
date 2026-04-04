package securestore

import (
	"strings"
	"testing"
)

func TestEncryptAndDecryptText(t *testing.T) {
	machineID := "test-machine-id"
	plainText := "sk-test-secret"

	encrypted, err := EncryptText(plainText, machineID)
	if err != nil {
		t.Fatalf("encrypt text failed: %v", err)
	}
	if encrypted == plainText {
		t.Fatal("expected encrypted value to differ from plain text")
	}
	if !strings.HasPrefix(encrypted, encryptedPrefix) {
		t.Fatalf("expected encrypted prefix %s", encryptedPrefix)
	}

	decrypted, err := DecryptText(encrypted, machineID)
	if err != nil {
		t.Fatalf("decrypt text failed: %v", err)
	}
	if decrypted != plainText {
		t.Fatalf("expected %s, got %s", plainText, decrypted)
	}
}

func TestDecryptTextWithDifferentMachineID(t *testing.T) {
	encrypted, err := EncryptText("sk-test-secret", "machine-a")
	if err != nil {
		t.Fatalf("encrypt text failed: %v", err)
	}

	if _, err = DecryptText(encrypted, "machine-b"); err == nil {
		t.Fatal("expected decrypt error with different machine id")
	}
}

func TestDecryptTextWithPlainTextInput(t *testing.T) {
	decrypted, err := DecryptText("plain-text", "test-machine-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decrypted != "plain-text" {
		t.Fatalf("expected plain-text, got %s", decrypted)
	}
}
