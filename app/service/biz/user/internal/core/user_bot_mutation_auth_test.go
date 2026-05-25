package core

import "testing"

func TestCanMutateBotByOwner(t *testing.T) {
	cases := []struct {
		name    string
		mdID    int64
		ownerID int64
		ok      bool
	}{
		{name: "match", mdID: 10, ownerID: 10, ok: true},
		{name: "md zero", mdID: 0, ownerID: 10, ok: false},
		{name: "owner zero", mdID: 10, ownerID: 0, ok: false},
		{name: "mismatch", mdID: 10, ownerID: 11, ok: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := canMutateBotByOwner(tc.mdID, tc.ownerID); got != tc.ok {
				t.Fatalf("canMutateBotByOwner(%d,%d)=%v want %v", tc.mdID, tc.ownerID, got, tc.ok)
			}
		})
	}
}
