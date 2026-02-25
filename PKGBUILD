# Maintainer: Atila de Freitas <atiladefreitas@gmail.com>
pkgname=prevy
pkgver=0.1.0
pkgrel=1
pkgdesc="A minimal terminal clipboard history manager built with Go and Bubble Tea"
arch=('x86_64' 'aarch64')
url="https://github.com/atiladefreitas/prevy"
license=('MIT')
depends=('glibc')
makedepends=('go')
optdepends=(
  'wl-clipboard: Wayland clipboard support'
  'xclip: X11 clipboard support'
  'xsel: X11 clipboard support (alternative)'
)
source=("$pkgname-$pkgver.tar.gz::https://github.com/atiladefreitas/prevy/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
  cd "$pkgname-$pkgver"
  export CGO_ENABLED=0
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  go build -ldflags "-s -w -X main.version=$pkgver" -o prevy .
}

package() {
  cd "$pkgname-$pkgver"
  install -Dm755 prevy "$pkgdir/usr/bin/prevy"
  install -Dm644 prevy.service "$pkgdir/usr/lib/systemd/user/prevy.service"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
  install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
