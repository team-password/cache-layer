package builder

import (
	"testing"
)

type MockEntry struct {
	RelateId   int64 `cache:"relateId"`
	SourceId   int64 `cache:"sourceId"`
	PropertyId int64 `cache:"propertyId"`
}

func (m MockEntry) TableName() string {
	return "public_relation"
}

type MockEntry2 struct {
	Id         int64
	RelateId   int64
	SourceId   int64
	PropertyId int64
}

func (m MockEntry2) TableName() string {
	return "public_relation"
}

func TestFmtSql(t *testing.T) {
	var a *int64
	var i int64 = 2
	var b = &i

	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id = ?",
				args: []interface{}{1},
			},
			want: "SELECT * FROM spu WHERE id = 1",
		},
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id in ?",
				args: []interface{}{[]string{"1", "2"}},
			},
			want: "SELECT * FROM spu WHERE id in (1,2)",
		},
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id =  ? and id = ? limit ?,?",
				args: []interface{}{1, "2", 0, 10},
			},
			want: "SELECT * FROM spu WHERE id =  1 and id = 2 limit 0,10",
		},
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id =  ? and id = ? limit ?,?",
				args: []interface{}{nil, a, b, []int{}},
			},
			want: "SELECT * FROM spu WHERE id =   ( 1 != 1 )  and id =  ( 1 != 1 )  limit 2, ( 1 != 1 ) ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSql(tt.args.sql, tt.args.args...); got != tt.want {
				t.Errorf("GenerateSql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateCountSql(t *testing.T) {
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id =  ?",
				args: []interface{}{"1"},
			},
			want: "SELECT COUNT(*) FROM (SELECT * FROM spu WHERE id =  1) t",
		},
		{
			name: "",
			args: args{
				sql:  "SELECT * FROM spu WHERE id in (1,2,3) LIMIT ?,?",
				args: []interface{}{0, 10},
			},
			want: "SELECT COUNT(*) FROM (SELECT * FROM spu WHERE id in (1,2,3) ) t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCountSql(tt.args.sql, tt.args.args...); got != tt.want {
				t.Errorf("GenerateCountSql() = %v, want %v", got, tt.want)
			}
		})
	}
}
