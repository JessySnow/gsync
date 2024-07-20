package index

import (
	"testing"
)

var (
	ar1 = RepoRelease{RepoOwner: "Tester", RepoName: "A", ReleaseName: "AR1", ReleaseTime: "2022-06-05T12:37:28Z"}
	ar2 = RepoRelease{RepoOwner: "Tester", RepoName: "A", ReleaseName: "AR2", ReleaseTime: "2022-06-05T12:37:29Z"}
	ar3 = RepoRelease{RepoOwner: "Tester", RepoName: "A", ReleaseName: "AR3", ReleaseTime: "2022-06-05T12:37:30Z"}
	ar4 = RepoRelease{RepoOwner: "Tester", RepoName: "A", ReleaseName: "AR4", ReleaseTime: "2022-06-05T12:38:30Z"}
	ar5 = RepoRelease{RepoOwner: "Tester", RepoName: "A", ReleaseName: "AR5", ReleaseTime: "2022-06-05T12:39:30Z"}
	br1 = RepoRelease{RepoOwner: "Tester", RepoName: "B", ReleaseName: "BR1", ReleaseTime: "2022-06-05T12:37:28Z"}
	br2 = RepoRelease{RepoOwner: "Tester", RepoName: "B", ReleaseName: "BR2", ReleaseTime: "2022-06-05T12:38:29Z"}
	br3 = RepoRelease{RepoOwner: "Tester", RepoName: "B", ReleaseName: "BR3", ReleaseTime: "2022-06-05T12:39:29Z"}
)

func TestBptreeReleaseDirIndexer_Add(t *testing.T) {
	indexer, err := New()
	if err != nil {
		t.Errorf("failed to init a indexer: %v", err)
		return
	}

	node, err := indexer.Add(&ar1)
	if err != nil {
		t.Errorf("failed to add reporelease: %v", err)
		return
	}
	dirName, _ := ar1.GenerateReleaseDirName()
	if node.DirName != dirName {
		t.Errorf("mismatch on reporelease and dirnode, expect: %s, got: %s", dirName, node.DirName)
		return
	}
}

func TestBptreeReleaseDirIndexer_GetAbsent(t *testing.T) {
	indexer, err := New()
	if err != nil {
		t.Errorf("failed to init a indexer: %v", err)
		return
	}

	// check in empty state
	absent, err := indexer.GetAbsent([]RepoRelease{ar1, ar2, br1, br2})
	if err != nil {
		t.Errorf("failed to get absent: %v", err)
		return
	}
	if len(absent) != 4 {
		t.Errorf("get absent mismatch: except absent reporelease: %d, got absent reporelease: %d", 4, len(absent))
		return
	}

	_, _ = indexer.Add(&ar3)
	_, _ = indexer.Add(&ar1)
	_, _ = indexer.Add(&br1)
	_, _ = indexer.Add(&br2)
	_, _ = indexer.Add(&br3)
	absent, err = indexer.GetAbsent([]RepoRelease{ar1, ar2, ar3, ar4, ar5, br1, br2, br3})
	if err != nil {
		t.Errorf("failed to get absent: %v", err)
		return
	}
	if len(absent) != 3 {
		t.Errorf("get absent mismatch: except absent reporelease: %d, got absent reporelease: %d", 3, len(absent))
		return
	}

	absentSet := make(map[RepoRelease]int)
	for _, abs := range absent {
		absentSet[abs] = 1
	}

	if absentSet[ar2] == 0 || absentSet[ar5] == 0 || absentSet[ar4] == 0 {
		t.Errorf("mismatch on except absent release and got absent release: %v, %v", []RepoRelease{ar2, br1, br2}, absentSet)
	}
}

func TestBptreeReleaseDirIndexer_Locate(t *testing.T) {
	indexer, err := New()
	if err != nil {
		t.Errorf("failed to init a indexer: %v", err)
		return
	}

	_, _ = indexer.Add(&ar3)
	_, _ = indexer.Add(&ar1)
	_, _ = indexer.Add(&br1)
	_, _ = indexer.Add(&br2)
	_, _ = indexer.Add(&br3)

	rrs := []RepoRelease{ar1, ar3, br1, br2, br3}
	for _, rr := range rrs {
		node, err := indexer.Locate(&rr)
		if err != nil {
			t.Errorf("failed to locate reporelease: %v", err)
			return
		}
		dirName, err := rr.GenerateReleaseDirName()
		if node.DirName != dirName {
			t.Errorf("mismatch on locate node, except: %v, got: %v", rr, node)
			return
		}
	}
}

func TestBptreeReleaseDirIndexer_Update(t *testing.T) {
	indexer, err := New()
	if err != nil {
		t.Errorf("failed to init a indexer: %v", err)
		return
	}

	_, _ = indexer.Add(&ar3)
	_, _ = indexer.Add(&ar1)
	_, _ = indexer.Add(&br1)
	_, _ = indexer.Add(&br2)
	_, _ = indexer.Add(&br3)

	node, _ := indexer.Locate(&br2)
	node.Locked = true
	err = indexer.Update(&br2, node)
	if err != nil {
		t.Errorf("update node failed: %v", err)
		return
	}
	updatedNode, _ := indexer.Locate(&br2)
	if !updatedNode.Locked {
		t.Errorf("mismatch on locked state, except: %v, got: %v", true, updatedNode.Locked)
		return
	}
}

func Test_generateKey(t *testing.T) {
	type args struct {
		rr *RepoRelease
	}
	tests := []struct {
		name    string
		args    args
		wantKey int
		wantErr bool
	}{
		{name: "ar1", wantErr: false, args: args{rr: &ar1}, wantKey: 1739129062},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := generateKey(tt.args.rr)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("generateKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
