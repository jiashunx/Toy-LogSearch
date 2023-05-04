package utils

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGetCommand(t *testing.T) {
    ast := assert.New(t)
    ast.Equal(CT1, GetCommand("sit2 newcore 2023-03-19 17\\:17\\:17 202304181703222101000066"))
    ast.Equal(CT2, GetCommand("sit2 newcore grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066"))
    ast.Equal(CT3, GetCommand("sit2 newcore bash cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066"))
}

func TestGetCommandArgs(t *testing.T) {
    ast := assert.New(t)

    args0 := GetCommandArgs("sit2 newcore ")
    ast.Empty(args0)

    args1 := GetCommandArgs("sit2 newcore 2023-03-19 17\\:17\\:17 202304181703222101000066")
    ast.Equal(5, len(args1))
    ast.Equal("sit2", args1[0])
    ast.Equal("newcore", args1[1])
    ast.Equal("2023-03-19", args1[2])
    ast.Equal("17\\:17\\:17", args1[3])
    ast.Equal("202304181703222101000066", args1[4])

    args2 := GetCommandArgs("sit2 newcore grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066")
    ast.Equal(4, len(args2));
    ast.Equal("sit2", args2[0]);
    ast.Equal("newcore", args2[1]);
    ast.Equal("grep", args2[2]);
    ast.Equal("2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066", args2[3]);

    args3 := GetCommandArgs("sit2 newcore bash cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066")
    ast.Equal(4, len(args3));
    ast.Equal("sit2", args3[0]);
    ast.Equal("newcore", args3[1]);
    ast.Equal("bash", args3[2]);
    ast.Equal("cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066", args3[3]);
}
