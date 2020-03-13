db.createUser(
    {
        user: "mongoadmin",
        pwd: "secret",
        roles: [
            {
                role: "readWrite",
                db: "platform"
            }
        ]
    }
);