package mysql

import "testing"

func TestGenerateSelectQuery(t *testing.T) {
	type args struct {
		selection  []string
		tableName  string
		conditions map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"GenerateSelectQuery", args{[]string{"id", "name", "email"}, "admins", map[string]interface{}{"id": 1}}, "SELECT id, name, email FROM admins WHERE id='1'"},
		{"GenerateSelectQuery", args{[]string{"id", "email"}, "admins", map[string]interface{}{"id": 1, "email": "mandip@gmail.com"}}, "SELECT id, email FROM admins WHERE id='1' AND email='mandip@gmail.com'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSelectQuery(tt.args.selection, tt.args.tableName, tt.args.conditions)
			if got != tt.want {
				t.Errorf("GenerateSelectQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
