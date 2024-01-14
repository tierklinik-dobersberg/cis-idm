# SMS (Twilio) Setup

`cisidm` supports phone number verifications and two-factor authentication using
SMS codes. For this to work you need to configure Twilio support in `cisidm` by
adding the `twilio` configuration block to your configuration file:

```hcl
# The (single) twilio block configures the Twilio integration which allows
# sending SMS messages to your users. This is required for phone-number
# verification to work.
twilio {
    from = "Example Inc"
    sid = "your-twilio-account-sid"
    token = "your-twilio-account-token"
}

```