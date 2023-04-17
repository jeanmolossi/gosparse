package gosparse_test

import (
	gosparse "gosparse"
	"testing"
)

func TestGetPredicate(t *testing.T) {
	type args struct {
		pre string
	}

	tests := []struct {
		name string
		args args
		want gosparse.Predicate
	}{
		{
			name: "valid field with predicate",
			args: args{pre: "fields[name_in]"},
			want: gosparse.IN,
		},
		{
			name: "valid field without predicate",
			args: args{pre: "fields[name]"},
			want: gosparse.NONE,
		},
		{
			name: "invalid field with predicate",
			args: args{pre: "fields[_nin]"},
			want: gosparse.NONE,
		},
		{
			name: "invalid field without",
			args: args{pre: "fields[]"},
			want: gosparse.NONE,
		},
		{
			name: "invalid field with predicate invalid",
			args: args{pre: "fields[_]"},
			want: gosparse.NONE,
		},
		// TEST ALL PREDICATES
		{
			name: "valid field with predicate 1",
			args: args{pre: "fields[name_eq]"},
			want: gosparse.EQ,
		},
		{
			name: "valid field with predicate 2",
			args: args{pre: "fields[name_neq]"},
			want: gosparse.NEQ,
		},
		{
			name: "valid field with predicate 3",
			args: args{pre: "fields[name_in]"},
			want: gosparse.IN,
		},
		{
			name: "valid field with predicate 4",
			args: args{pre: "fields[name_nin]"},
			want: gosparse.NIN,
		},
		{
			name: "valid field with predicate 5",
			args: args{pre: "fields[name_gt]"},
			want: gosparse.GT,
		},
		{
			name: "valid field with predicate 6",
			args: args{pre: "fields[name_gte]"},
			want: gosparse.GTE,
		},
		{
			name: "valid field with predicate 7",
			args: args{pre: "fields[name_lt]"},
			want: gosparse.LT,
		},
		{
			name: "valid field with predicate 8",
			args: args{pre: "fields[name_lte]"},
			want: gosparse.LTE,
		},
		{
			name: "valid field with predicate 9",
			args: args{pre: "fields[name_blank]"},
			want: gosparse.BLANK,
		},
		{
			name: "valid field with predicate 10",
			args: args{pre: "fields[name_null]"},
			want: gosparse.NULL,
		},
		{
			name: "valid field with predicate 11",
			args: args{pre: "fields[name_notnull]"},
			want: gosparse.NOT_NULL,
		},
		{
			name: "valid field with predicate 12",
			args: args{pre: "fields[name_start]"},
			want: gosparse.START,
		},
		{
			name: "valid field with predicate 13",
			args: args{pre: "fields[name_end]"},
			want: gosparse.END,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gosparse.GetPredicate(tt.args.pre); got != tt.want {
				t.Errorf("GetPredicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
