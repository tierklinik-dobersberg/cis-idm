package cisidm.forward_auth

import future.keywords.in

default allow = false

headers["Foo"] := "Bar"

allow {
    user_is_root
}

user_is_root {
    input.subject
    input.subject.roles

    some role in input.subject.roles
    role = "idm_superuser"
}