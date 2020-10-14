github.com/teamlint/baron/cmd/baron
  ├ fmt
  ├ go/build
  ├ io
  ├ os
  ├ os/exec
  ├ path/filepath
  ├ strings
  ├ github.com/teamlint/baron/config
  │ └ io
  ├ github.com/teamlint/baron/gengokit
  │ ├ bytes
  │ ├ io
  │ ├ strings
  │ ├ text/template
  │ ├ github.com/teamlint/baron/gengokit/httptransport
  │ │ ├ bytes
  │ │ ├ fmt
  │ │ ├ go/ast
  │ │ ├ go/format
  │ │ ├ go/parser
  │ │ ├ go/printer
  │ │ ├ go/token
  │ │ ├ reflect
  │ │ ├ runtime
  │ │ ├ strconv
  │ │ ├ strings
  │ │ ├ text/template
  │ │ ├ unicode
  │ │ ├ github.com/teamlint/baron/gengokit/httptransport/templates
  │ │ ├ github.com/teamlint/baron/svcdef
  │ │ │ ├ fmt
  │ │ │ ├ go/ast
  │ │ │ ├ go/parser
  │ │ │ ├ go/token
  │ │ │ ├ io
  │ │ │ ├ io/ioutil
  │ │ │ ├ os
  │ │ │ ├ path/filepath
  │ │ │ ├ reflect
  │ │ │ ├ regexp
  │ │ │ ├ strings
  │ │ │ ├ github.com/teamlint/baron/svcdef/svcparse
  │ │ │ │ ├ bufio
  │ │ │ │ ├ bytes
  │ │ │ │ ├ fmt
  │ │ │ │ ├ io
  │ │ │ │ ├ io/ioutil
  │ │ │ │ ├ strconv
  │ │ │ │ ├ strings
  │ │ │ │ ├ unicode
  │ │ │ │ └ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ │ ├ github.com/teamlint/baron/baron/execprotoc
  │ │ │ │ ├ io/ioutil
  │ │ │ │ ├ os
  │ │ │ │ ├ os/exec
  │ │ │ │ ├ path/filepath
  │ │ │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/proto
  │ │ │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/plugin
  │ │ │ │ └ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ │ └ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator
  │ │ ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ └ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  │ ├ github.com/teamlint/baron/svcdef
  │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator
  │ │ ├ bufio
  │ │ ├ bytes
  │ │ ├ compress/gzip
  │ │ ├ crypto/sha256
  │ │ ├ encoding/hex
  │ │ ├ fmt
  │ │ ├ go/ast
  │ │ ├ go/build
  │ │ ├ go/parser
  │ │ ├ go/printer
  │ │ ├ go/token
  │ │ ├ log
  │ │ ├ os
  │ │ ├ path
  │ │ ├ sort
  │ │ ├ strconv
  │ │ ├ strings
  │ │ ├ unicode
  │ │ ├ unicode/utf8
  │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/gogoproto
  │ │ │ ├ fmt
  │ │ │ ├ math
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/proto
  │ │ │ │ ├ bufio
  │ │ │ │ ├ bytes
  │ │ │ │ ├ encoding
  │ │ │ │ ├ encoding/json
  │ │ │ │ ├ errors
  │ │ │ │ ├ fmt
  │ │ │ │ ├ io
  │ │ │ │ ├ log
  │ │ │ │ ├ math
  │ │ │ │ ├ reflect
  │ │ │ │ ├ sort
  │ │ │ │ ├ strconv
  │ │ │ │ ├ strings
  │ │ │ │ ├ sync
  │ │ │ │ ├ sync/atomic
  │ │ │ │ ├ time
  │ │ │ │ ├ unicode/utf8
  │ │ │ │ └ unsafe
  │ │ │ └ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/descriptor
  │ │ │   ├ bytes
  │ │ │   ├ compress/gzip
  │ │ │   ├ fmt
  │ │ │   ├ io/ioutil
  │ │ │   ├ math
  │ │ │   ├ reflect
  │ │ │   ├ sort
  │ │ │   ├ strconv
  │ │ │   ├ strings
  │ │ │   └ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/proto
  │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/proto
  │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/descriptor
  │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator/internal/remap
  │ │ │ ├ fmt
  │ │ │ ├ go/scanner
  │ │ │ └ go/token
  │ │ └ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/plugin
  │ │   ├ fmt
  │ │   ├ math
  │ │   ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/proto
  │ │   └ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/descriptor
  │ └ github.com/teamlint/baron/vendor/github.com/pkg/errors
  ├ github.com/teamlint/baron/gengokit/generator
  │ ├ bytes
  │ ├ go/format
  │ ├ io
  │ ├ io/ioutil
  │ ├ strings
  │ ├ github.com/teamlint/baron/gengokit
  │ ├ github.com/teamlint/baron/gengokit/handlers
  │ │ ├ bytes
  │ │ ├ go/ast
  │ │ ├ go/parser
  │ │ ├ go/printer
  │ │ ├ go/token
  │ │ ├ io
  │ │ ├ strings
  │ │ ├ github.com/teamlint/baron/gengokit
  │ │ │ ├ bytes
  │ │ │ ├ io
  │ │ │ ├ strings
  │ │ │ ├ text/template
  │ │ │ ├ github.com/teamlint/baron/gengokit/httptransport
  │ │ │ ├ github.com/teamlint/baron/svcdef
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator
  │ │ │ └ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ ├ github.com/teamlint/baron/gengokit/handlers/templates
  │ │ ├ github.com/teamlint/baron/gengokit
  │ │ ├ github.com/teamlint/baron/gengokit/handlers/templates
  │ │ ├ github.com/teamlint/baron/svcdef
  │ │ │ ├ fmt
  │ │ │ ├ go/ast
  │ │ │ ├ go/parser
  │ │ │ ├ go/token
  │ │ │ ├ io
  │ │ │ ├ io/ioutil
  │ │ │ ├ os
  │ │ │ ├ path/filepath
  │ │ │ ├ reflect
  │ │ │ ├ regexp
  │ │ │ ├ strings
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ │ ├ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  │ │ │ ├ github.com/teamlint/baron/svcdef/svcparse (unresolved)
  │ │ │ └ github.com/teamlint/baron/baron/execprotoc (unresolved)
  │ │ ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ │ └ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  │ ├ github.com/teamlint/baron/gengokit/template
  │ │ ├ bytes
  │ │ ├ compress/gzip
  │ │ ├ crypto/sha256
  │ │ ├ fmt
  │ │ ├ io
  │ │ ├ io/ioutil
  │ │ ├ os
  │ │ ├ path/filepath
  │ │ ├ strings
  │ │ └ time
  │ ├ github.com/teamlint/baron/svcdef
  │ ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ └ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  ├ github.com/teamlint/baron/internal/execprotoc
  │ ├ io/ioutil
  │ ├ os
  │ ├ os/exec
  │ ├ path/filepath
  │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/proto
  │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/plugin
  │ └ github.com/teamlint/baron/vendor/github.com/pkg/errors
  ├ github.com/teamlint/baron/internal/parsesvcname
  │ ├ io
  │ ├ io/ioutil
  │ ├ os
  │ ├ path/filepath
  │ ├ strings
  │ ├ github.com/teamlint/baron/internal/execprotoc
  │ ├ github.com/teamlint/baron/svcdef
  │ └ github.com/teamlint/baron/vendor/github.com/pkg/errors
  ├ github.com/teamlint/baron/internal/start
  │ ├ bytes
  │ ├ os
  │ ├ strings
  │ ├ text/template
  │ ├ github.com/teamlint/baron/vendor/github.com/gogo/protobuf/protoc-gen-gogo/generator
  │ ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ └ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  ├ github.com/teamlint/baron/svcdef
  ├ github.com/teamlint/baron/vendor/github.com/pkg/errors
  │ ├ fmt
  │ ├ io
  │ ├ path
  │ ├ runtime
  │ └ strings
  ├ github.com/teamlint/baron/vendor/github.com/sirupsen/logrus
  │ ├ bufio
  │ ├ bytes
  │ ├ context
  │ ├ encoding/json
  │ ├ fmt
  │ ├ io
  │ ├ log
  │ ├ os
  │ ├ reflect
  │ ├ runtime
  │ ├ sort
  │ ├ strings
  │ ├ sync
  │ ├ sync/atomic
  │ ├ time
  │ └ github.com/teamlint/baron/vendor/golang.org/x/sys/unix
  │   ├ bytes
  │   ├ errors
  │   ├ runtime
  │   ├ sort
  │   ├ strings
  │   ├ sync
  │   ├ syscall
  │   ├ time
  │   └ unsafe
  ├ github.com/teamlint/baron/vendor/github.com/spf13/pflag
  │ ├ bytes
  │ ├ encoding/base64
  │ ├ encoding/csv
  │ ├ encoding/hex
  │ ├ errors
  │ ├ flag
  │ ├ fmt
  │ ├ io
  │ ├ net
  │ ├ os
  │ ├ reflect
  │ ├ sort
  │ ├ strconv
  │ ├ strings
  │ └ time
  └ github.com/teamlint/baron/vendor/golang.org/x/tools/go/packages
    ├ bytes
    ├ context
    ├ encoding/json
    ├ fmt
    ├ go/ast
    ├ go/parser
    ├ go/scanner
    ├ go/token
    ├ go/types
    ├ io/ioutil
    ├ log
    ├ os
    ├ os/exec
    ├ path
    ├ path/filepath
    ├ reflect
    ├ regexp
    ├ sort
    ├ strconv
    ├ strings
    ├ sync
    ├ time
    ├ unicode
    ├ github.com/teamlint/baron/vendor/golang.org/x/tools/go/gcexportdata
    │ ├ bufio
    │ ├ bytes
    │ ├ fmt
    │ ├ go/token
    │ ├ go/types
    │ ├ io
    │ ├ io/ioutil
    │ ├ os
    │ └ github.com/teamlint/baron/vendor/golang.org/x/tools/go/internal/gcimporter
    │   ├ bufio
    │   ├ bytes
    │   ├ encoding/binary
    │   ├ errors
    │   ├ fmt
    │   ├ go/ast
    │   ├ go/build
    │   ├ go/constant
    │   ├ go/token
    │   ├ go/types
    │   ├ io
    │   ├ io/ioutil
    │   ├ math
    │   ├ math/big
    │   ├ os
    │   ├ path/filepath
    │   ├ reflect
    │   ├ sort
    │   ├ strconv
    │   ├ strings
    │   ├ sync
    │   ├ text/scanner
    │   ├ unicode
    │   └ unicode/utf8
    ├ github.com/teamlint/baron/vendor/golang.org/x/tools/go/internal/packagesdriver
    │ ├ bytes
    │ ├ context
    │ ├ encoding/json
    │ ├ fmt
    │ ├ go/types
    │ ├ log
    │ ├ os
    │ ├ os/exec
    │ ├ strings
    │ └ time
    ├ github.com/teamlint/baron/vendor/golang.org/x/tools/internal/gopathwalk
    │ ├ bufio
    │ ├ bytes
    │ ├ fmt
    │ ├ go/build
    │ ├ io/ioutil
    │ ├ log
    │ ├ os
    │ ├ path/filepath
    │ ├ strings
    │ └ github.com/teamlint/baron/vendor/golang.org/x/tools/internal/fastwalk
    │   ├ errors
    │   ├ fmt
    │   ├ os
    │   ├ path/filepath
    │   ├ runtime
    │   ├ sync
    │   ├ syscall
    │   └ unsafe
    ├ github.com/teamlint/baron/vendor/golang.org/x/tools/internal/semver
    └ github.com/teamlint/baron/vendor/golang.org/x/tools/internal/span
      ├ encoding/json
      ├ fmt
      ├ go/token
      ├ net/url
      ├ os
      ├ path
      ├ path/filepath
      ├ runtime
      ├ strconv
      ├ strings
      ├ unicode
      ├ unicode/utf16
      └ unicode/utf8
91 dependencies (50 internal, 41 external, 0 testing).
