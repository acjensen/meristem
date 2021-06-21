# meristem

Plant growth animation represented as a simple discrete [L-system](https://en.wikipedia.org/wiki/L-system).

For more info see my [blog post](http://www.acjensen.com/l-system/).

# Example
`go run meristem.go`

`convert $(ls -1 img/*.png | sort -V) -loop 0 animation.gif`
