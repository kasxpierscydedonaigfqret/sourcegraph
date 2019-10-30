	"reflect"
	"github.com/sourcegraph/sourcegraph/internal/vcs/git/gittest"
	t.Parallel()
	srv := httptest.NewServer((&server.Server{}).Handler())
	repoWithDotGitDir := gittest.MakeTmpDir(t, "repo-with-dot-git-dir")
	if err := createRepoWithDotGitDir(repoWithDotGitDir); err != nil {
		t.Fatal(err)
	}

	gitCommands := []string{
		"mkdir dir1",
		"echo -n infile1 > dir1/file1",
		"touch --date=2006-01-02T15:04:05Z dir1 dir1/file1 || touch -t " + gittest.Times[0] + " dir1 dir1/file1",
		"git add dir1/file1",
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit -m commit1 --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
		"echo -n infile2 > 'file 2'",
		"touch --date=2014-05-06T19:20:21Z 'file 2' || touch -t " + gittest.Times[1] + " 'file 2'",
		"git add 'file 2'",
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2014-05-06T19:20:21Z git commit -m commit2 --author='a <a@a.com>' --date 2014-05-06T19:20:21Z",
	}
	tests := map[string]struct {
		repo gitserver.Repo
		want map[string]string
		err  error
		"git cmd": {
			repo: gittest.MakeGitRepository(t, gitCommands...),
		"repo with .git dir": {
			repo: gitserver.Repo{Name: api.RepoName(repoWithDotGitDir), URL: repoWithDotGitDir},
			want: map[string]string{"file1": "hello\n", ".git/mydir/file2": "milton\n", ".git/mydir/": "", ".git/": ""},
		"repo not found": {
			repo: gitserver.Repo{Name: api.RepoName("not-found")},
			err:  errors.New("repository does not exist: not-found"),
	for label, test := range tests {
		rc, err := cli.Archive(ctx, test.repo, gitserver.ArchiveOptions{Treeish: "HEAD", Format: "zip"})
		if have, want := fmt.Sprint(err), fmt.Sprint(test.err); have != want {
			t.Errorf("%s: Archive: have err %v, want %v", label, have, want)
		}
		if rc == nil {
			continue
		}
		defer rc.Close()
		data, err := ioutil.ReadAll(rc)
		if err != nil {
			t.Errorf("%s: ReadAll: %s", label, err)
			continue
		}
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			t.Errorf("%s: zip.NewReader: %s", label, err)
			continue
		}
		got := map[string]string{}
		for _, f := range zr.File {
			r, err := f.Open()
				t.Errorf("%s: failed to open %q because %s", label, f.Name, err)
				continue
			contents, err := ioutil.ReadAll(r)
			r.Close()
				t.Errorf("%s: Read(%q): %s", label, f.Name, err)
				continue
			got[f.Name] = string(contents)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: got %v, want %v", label, got, test.want)
		}
func createRepoWithDotGitDir(dir string) error {
			panic(err)
			return err
			return err
	return nil