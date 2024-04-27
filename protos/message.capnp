using Go = import "/go.capnp";
@0xe945d32308a30635;
$Go.package("protos");
$Go.import("protos/message");

struct KoboldMessage $Go.doc("standard kobold message for transfering info between clients and nodes"){
    # identifier of the message
    id @0 :Text;

    # target topic for the message to be sent
    topic @1 :Text;

    # Transaction ID
    txId @2 :Text;

    # Content
    content @3 :Data;
}
