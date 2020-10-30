package start

const startProto = `
syntax = "proto3";

package {{.PackageName}};

import "google/api/annotations.proto";

service {{.ServiceName}} {
  rpc Status(StatusRequest) returns (StatusResponse) {
    option (google.api.http) = {
      get: "/status"
    };
  }
}

enum ServiceStatus {
  FAIL = 0;
  OK = 1;
}

message StatusRequest {
  bool full = 1;
}

message StatusResponse {
  ServiceStatus status = 1;
}
`

const nextStepMsg = `A "start" protobuf file named '{{.FileName}}' has been created in the
current directory. You can generate a service based on this new protobuf file
at any time using the following command:

    baron {{.FileName}}

If you want to generate a protofile with a different name, use the
'--start' option with the name of your choice after '--start'. For
example, to generate a 'foo.proto', use the following command:

    baron --start foo
`
const existingFileMsg = `There's already a "start" protobuf file named '{{.FileName}}' in the current
directory. If you'd like to generate a service based on this existing protobuf
file, you should instead run the command:

    baron {{.FileName}}`
