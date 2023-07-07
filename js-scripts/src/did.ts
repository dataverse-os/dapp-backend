import { DID } from "dids";
import { Ed25519Provider } from "key-did-provider-ed25519";
import { getResolver } from "key-did-resolver";
import { fromString } from "uint8arrays/from-string";

(async () => {
  const privateKey = fromString(process.argv[2], "base16");
  const did = new DID({
    resolver: getResolver(),
    provider: new Ed25519Provider(privateKey),
  });
  await did.authenticate();
  console.log(did.id);
})();