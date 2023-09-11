extern crate libc;
use anyhow::Result;
use async_std::task;
use multibase::Base;
use ssh_key::private::Ed25519Keypair;
use std::ffi::{CStr, CString};

use ceramic_http_client::{
    ceramic_event::{DidDocument, JwkSigner},
    remote::CeramicRemoteHttpClient,
};

fn generate_did_str(pk: &str) -> Result<String> {
    let seed: [u8; 32] = hex::decode(pk)?
        .try_into()
        .expect("seed length is 32 bytes");
    let key = Ed25519Keypair::from_seed(&seed);

    let mut buf: Vec<u8> = vec![0xed, 0x01];
    buf.extend(key.public.0);

    Ok(format!(
        "did:key:{}",
        multibase::encode(Base::Base58Btc, buf)
    ))
}

fn ceramic_client(ceramic: &str, pk: &str) -> Result<CeramicRemoteHttpClient<JwkSigner>> {
    let did = generate_did_str(pk)?;
    let did = DidDocument::new(&did);
    let signer = task::block_on(JwkSigner::new(did, pk))?;

    let ceramic_url = url::Url::parse(ceramic)?;
    Ok(CeramicRemoteHttpClient::new(signer, ceramic_url))
}

#[repr(C)]
pub struct CResult {
    data: *const libc::c_char,
    err: libc::c_int,
}

#[no_mangle]
pub extern "C" fn generate_did(key: *const libc::c_char) -> CResult {
    let pk = unsafe { CStr::from_ptr(key) }.to_str().unwrap();
    match generate_did_str(pk) {
        Ok(res) => CResult {
            data: CString::new(res).unwrap().into_raw(),
            err: 0 as libc::c_int,
        },
        Err(err) => CResult {
            data: CString::new(err.to_string()).unwrap().into_raw(),
            err: -1 as libc::c_int,
        },
    }
}

#[no_mangle]
pub extern "C" fn get_ceramic_node_status(
    ceramic: *const libc::c_char,
    key: *const libc::c_char,
) -> CResult {
    let pk = unsafe { CStr::from_ptr(key) }.to_str().unwrap();
    let ceramic = unsafe { CStr::from_ptr(ceramic) }.to_str().unwrap();

    let ceramic = match ceramic_client(ceramic, pk) {
        Ok(ceramic) => ceramic,
        Err(err) => {
            return CResult {
                data: CString::new(err.to_string()).unwrap().into_raw(),
                err: -1 as libc::c_int,
            }
        }
    };

    let rt = tokio::runtime::Runtime::new().unwrap();
    let res = match rt.block_on(ceramic.node_status()) {
        Ok(_) => "success",
        Err(_) => "failed",
    };
    CResult {
        data: CString::new(res).unwrap().into_raw(),
        err: 0 as libc::c_int,
    }
}
