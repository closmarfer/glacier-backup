//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -build_constraint=mocks

package backup

type Application interface {
	Run()
}
