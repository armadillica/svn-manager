package httphandler

import (
	"github.com/armadillica/svn-manager/svnman"
	check "gopkg.in/check.v1"
)

var validModifyGrantAccess = svnman.ModifyAccess{
	Grant: []svnman.ModifyAccessGrantEntry{
		svnman.ModifyAccessGrantEntry{
			Username: "joey",
			Password: "$2y$05$cWVQLHS58K7fIKjz3tU52eBI2sxbE3KdAfZN0CJN/DcRKGkYTKOuG", // jemoeder
		},
		svnman.ModifyAccessGrantEntry{
			Username: "strongman",
			Password: "$2y$05$YpqSmhP7x06Z05bkfnXlXu3z88mFzbIoVH5kY/p1eFQ2qC17BeyxG", // überstrong
		},
	},
}

var validModifyRevokeAccess = svnman.ModifyAccess{
	Revoke: []string{"joey"},
}

func (s *ValidationTestSuite) TestEmptyRequestHappy(t *check.C) {
	s.assertValidJSON(t, "modify_access", svnman.ModifyAccess{})
	s.assertValidJSON(t, "modify_access", svnman.ModifyAccess{
		Grant: []svnman.ModifyAccessGrantEntry{},
	})
	s.assertValidJSON(t, "modify_access", svnman.ModifyAccess{
		Revoke: []string{},
	})
	s.assertValidJSON(t, "modify_access", svnman.ModifyAccess{
		Grant:  []svnman.ModifyAccessGrantEntry{},
		Revoke: []string{},
	})
}

func (s *ValidationTestSuite) TestGrantHappy(t *check.C) {
	s.assertValidJSON(t, "modify_access", validModifyGrantAccess)
}

func (s *ValidationTestSuite) TestGrantUnhappy(t *check.C) {
	s.assertInvalidJSON(t, "modify_access", svnman.ModifyAccess{
		Grant: []svnman.ModifyAccessGrantEntry{
			svnman.ModifyAccessGrantEntry{
				Username: "invalid username",
				Password: "$2y$05$cWVQLHS58K7fIKjz3tU52eBI2sxbE3KdAfZN0CJN/DcRKGkYTKOuG",
			},
		},
	})
	s.assertInvalidJSON(t, "modify_access", svnman.ModifyAccess{
		Grant: []svnman.ModifyAccessGrantEntry{
			svnman.ModifyAccessGrantEntry{
				Username: "joey",
				Password: "unencrypted password",
			},
		},
	})
}

func (s *ValidationTestSuite) TestRevokeHappy(t *check.C) {
	s.assertValidJSON(t, "modify_access", validModifyRevokeAccess)
}

func (s *ValidationTestSuite) TestRevokeUnhappy(t *check.C) {
	s.assertInvalidJSON(t, "modify_access", svnman.ModifyAccess{Revoke: []string{"invalid username"}})
	s.assertInvalidJSON(t, "modify_access", svnman.ModifyAccess{Revoke: []string{"joey\nstrongman"}})
	s.assertInvalidJSON(t, "modify_access", svnman.ModifyAccess{Revoke: []string{"üsername"}})
}
