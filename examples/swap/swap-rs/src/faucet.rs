use reqwest::header::CONTENT_TYPE;
use reqwest::Client;
use serde::Deserialize;
use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct FaucetError(String);

impl fmt::Display for FaucetError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl Error for FaucetError {}

#[derive(Deserialize)]
struct FaucetResponse {
    _task: Option<String>,
    error: Option<String>,
}

pub async fn request_fund_from_faucet(
    address: &str,
    faucet_url: &str,
) -> Result<(), Box<dyn Error>> {
    let param_json = format!(r#"{{"FixedAmountRequest":{{"recipient":"{}"}}}}"#, address);

    let client = Client::new();
    let res = client
        .post(faucet_url)
        .header(CONTENT_TYPE, "application/json")
        .body(param_json)
        .send()
        .await
        .unwrap();

    if !res.status().is_success() {
        return Err(Box::new(FaucetError(format!(
            "post {} response code: {}",
            faucet_url,
            res.status()
        ))));
    }

    let response: FaucetResponse = res.json().await.unwrap();
    if let Some(error) = response.error {
        if !error.trim().is_empty() {
            return Err(Box::new(FaucetError(error)));
        }
    }

    Ok(())
}
