package sui_test

import (
	"reflect"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/stretchr/testify/require"
)

func TestNewResourceType(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    *sui.ResourceType
		wantErr bool
	}{
		{
			name: "no subtype",
			str:  "0x23::coin::Xxxx",
			want: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x23"),
				Module:     "coin",
				ObjectName: "Xxxx",
				SubTypes:   nil,
			},
		},
		{
			name: "one subtype",
			str:  "0x1::m1::f1<0x2::m2::f2>",
			want: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x1"),
				Module:     "m1",
				ObjectName: "f1",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x2"),
						Module:     "m2",
						ObjectName: "f2",
						SubTypes:   nil,
					},
				},
			},
		},
		{
			name: "multiple subtypes",
			str:  "0x111::aaa::AAA<0x222::bbb::BBB, 0x333::ccc::CCC, 0x444::ddd::DDD>",
			want: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x111"),
				Module:     "aaa",
				ObjectName: "AAA",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x222"),
						Module:     "bbb",
						ObjectName: "BBB",
						SubTypes:   nil,
					},
					{
						Address:    sui.MustAddressFromHex("0x333"),
						Module:     "ccc",
						ObjectName: "CCC",
						SubTypes:   nil,
					},
					{
						Address:    sui.MustAddressFromHex("0x444"),
						Module:     "ddd",
						ObjectName: "DDD",
						SubTypes:   nil,
					},
				},
			},
		},
		{
			name: "nested generics two levels",
			str:  "0x2::dynamic_field::Field<0x1::ascii::String, 0x2::balance::Balance<0x2::iota::IOTA>>",
			want: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x2"),
				Module:     "dynamic_field",
				ObjectName: "Field",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x1"),
						Module:     "ascii",
						ObjectName: "String",
						SubTypes:   nil,
					},
					{
						Address:    sui.MustAddressFromHex("0x2"),
						Module:     "balance",
						ObjectName: "Balance",
						SubTypes: []*sui.ResourceType{
							{
								Address:    sui.MustAddressFromHex("0x2"),
								Module:     "iota",
								ObjectName: "IOTA",
								SubTypes:   nil,
							},
						},
					},
				},
			},
		},
		{
			name: "deeply nested generics",
			str:  "0x1::outer::O<0x2::mid::M<0x3::inner::I, 0x4::k::K>, 0x5::leaf::L>",
			want: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x1"),
				Module:     "outer",
				ObjectName: "O",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x2"),
						Module:     "mid",
						ObjectName: "M",
						SubTypes: []*sui.ResourceType{
							{
								Address:    sui.MustAddressFromHex("0x3"),
								Module:     "inner",
								ObjectName: "I",
								SubTypes:   nil,
							},
							{
								Address:    sui.MustAddressFromHex("0x4"),
								Module:     "k",
								ObjectName: "K",
								SubTypes:   nil,
							},
						},
					},
					{
						Address:    sui.MustAddressFromHex("0x5"),
						Module:     "leaf",
						ObjectName: "L",
						SubTypes:   nil,
					},
				},
			},
		},
		{
			name: "whitespace tolerance",
			str:  "0x1::m1::f1< 0x2::m2::f2 , 0x3::m3::f3 >",
			want: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x1"),
				Module:     "m1",
				ObjectName: "f1",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x2"),
						Module:     "m2",
						ObjectName: "f2",
						SubTypes:   nil,
					},
					{
						Address:    sui.MustAddressFromHex("0x3"),
						Module:     "m3",
						ObjectName: "f3",
						SubTypes:   nil,
					},
				},
			},
		},
		{
			name:    "error address",
			str:     "0x123abcg::coin::Xxxx",
			wantErr: true,
		},
		{
			name:    "error format - trailing chars",
			str:     "0x1::m1::f1<0x2::m2::f2>x",
			wantErr: true,
		},
		{
			name:    "error format - mismatched brackets",
			str:     "0x1::m1::f1<0x2::m2::f2<0x3::m3::f3>",
			wantErr: true,
		},
		{
			name:    "error format - leading bracket",
			str:     "<0x3::m3::f3>0x1::m1::f1<0x2::m2::f2>",
			wantErr: true,
		},
		{
			name:    "error format - empty subtype",
			str:     "0x1::m1::f1<0x2::m2::f2,>",
			wantErr: true,
		},
		{
			name:    "error format - empty generic",
			str:     "0x1::m1::f1<>",
			wantErr: true,
		},
		{
			name:    "error format - wrong parts count",
			str:     "0x1::m1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := sui.NewResourceType(tt.str)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewResourceType() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewResourceType(): %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		target *sui.ResourceType
		want   bool
	}{
		{
			name:   "successful, matches outer type",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Module: "swap", ObjectName: "Pool"},
			want:   true,
		},
		{
			name:   "successful, matches single subtype",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Module: "testcoin", ObjectName: "TESTCOIN"},
			want:   true,
		},
		{
			name:   "successful, matches first subtype in multiple subtypes",
			str:    "0x2::dynamic_field::Field<0x1::ascii::String, 0x2::balance::Balance<0x2::iota::IOTA>>",
			target: &sui.ResourceType{Module: "ascii", ObjectName: "String"},
			want:   true,
		},
		{
			name:   "successful, matches second subtype in multiple subtypes",
			str:    "0x2::dynamic_field::Field<0x1::ascii::String, 0x2::balance::Balance<0x2::iota::IOTA>>",
			target: &sui.ResourceType{Module: "balance", ObjectName: "Balance"},
			want:   true,
		},
		{
			name:   "successful, matches deeply nested subtype",
			str:    "0x2::dynamic_field::Field<0x1::ascii::String, 0x2::balance::Balance<0x2::iota::IOTA>>",
			target: &sui.ResourceType{Module: "iota", ObjectName: "IOTA"},
			want:   true,
		},
		{
			name:   "successful, matches in complex nested structure",
			str:    "0x111::aaa::AAA<0x222::bbb::BBB, 0x333::ccc::CCC, 0x444::ddd::DDD>",
			target: &sui.ResourceType{Module: "ccc", ObjectName: "CCC"},
			want:   true,
		},
		{
			name:   "successful, matches deeply nested type",
			str:    "0x1::outer::O<0x2::mid::M<0x3::inner::I, 0x4::k::K>, 0x5::leaf::L>",
			target: &sui.ResourceType{Module: "inner", ObjectName: "I"},
			want:   true,
		},
		{
			name:   "successful, with address match",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Address: sui.MustAddressFromHex("0xe87e"), Module: "swap", ObjectName: "Pool"},
			want:   true,
		},
		{
			name:   "failed, wrong module name",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Module: "name", ObjectName: "Pool"},
			want:   false,
		},
		{
			name:   "failed, wrong address",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Address: sui.MustAddressFromHex("0x1"), Module: "swap", ObjectName: "Pool"},
			want:   false,
		},
		{
			name:   "failed, not found anywhere",
			str:    "0x111::aaa::AAA<0x222::bbb::BBB, 0x333::ccc::CCC>",
			target: &sui.ResourceType{Module: "nonexistent", ObjectName: "NotFound"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				src, err := sui.NewResourceType(tt.str)
				require.NoError(t, err)
				require.Equal(t, tt.want, src.Contains(tt.target.Address, tt.target.Module, tt.target.ObjectName))
			},
		)
	}
}

func TestResourceTypeString(t *testing.T) {
	typeString := "0x1::mmm1::fff1<0x123abcdef::mm2::ff3>"

	resourceType, err := sui.NewResourceType(typeString)
	require.NoError(t, err)
	res := "0x0000000000000000000000000000000000000000000000000000000000000001::mmm1::fff1<0x0000000000000000000000000000000000000000000000000000000123abcdef::mm2::ff3>"
	require.Equal(t, resourceType.String(), res)
}

func TestResourceTypeNormalizationAndEquality(t *testing.T) {
	tests := []struct {
		name string
		str1 string
		str2 string
		want bool
	}{
		{
			name: "equivalent normalized addresses",
			str1: "0x1::module::Name<0x2::sub::Type>",
			str2: "0x0000000000000000000000000000000000000000000000000000000000000001::module::Name<0x0000000000000000000000000000000000000000000000000000000000000002::sub::Type>",
			want: true,
		},
		{
			name: "different addresses",
			str1: "0x1::module::Name",
			str2: "0x2::module::Name",
			want: false,
		},
		{
			name: "different modules",
			str1: "0x1::module1::Name",
			str2: "0x1::module2::Name",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			same, err := sui.IsSameResource(tt.str1, tt.str2)
			require.NoError(t, err)
			require.Equal(t, tt.want, same)
		})
	}
}

func TestResourceTypeShortString(t *testing.T) {
	tests := []struct {
		name string
		arg  *sui.ResourceType
		want string
	}{
		{
			name: "no subtypes",
			arg: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x1"),
				Module:     "m1",
				ObjectName: "f1",
				SubTypes:   nil,
			},
			want: "0x1::m1::f1",
		},
		{
			name: "nested generics",
			arg: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x1"),
				Module:     "m1",
				ObjectName: "f1",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x2"),
						Module:     "m2",
						ObjectName: "f2",
						SubTypes: []*sui.ResourceType{
							{
								Address:    sui.MustAddressFromHex("0x123abcdef"),
								Module:     "m3",
								ObjectName: "f3",
								SubTypes:   nil,
							},
						},
					},
				},
			},
			want: "0x1::m1::f1<0x2::m2::f2<0x123abcdef::m3::f3>>",
		},
		{
			name: "multiple subtypes",
			arg: &sui.ResourceType{
				Address:    sui.MustAddressFromHex("0x1"),
				Module:     "outer",
				ObjectName: "Type",
				SubTypes: []*sui.ResourceType{
					{
						Address:    sui.MustAddressFromHex("0x2"),
						Module:     "mod1",
						ObjectName: "Type1",
						SubTypes:   nil,
					},
					{
						Address:    sui.MustAddressFromHex("0x3"),
						Module:     "mod2",
						ObjectName: "Type2",
						SubTypes:   nil,
					},
				},
			},
			want: "0x1::outer::Type<0x2::mod1::Type1, 0x3::mod2::Type2>",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.arg.ShortString(); got != tt.want {
					t.Errorf("ResourceType.ShortString(): %v, want %v", got, tt.want)
				}
			},
		)
	}
}
