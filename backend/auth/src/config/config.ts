// storing varibles centrally
export const config = {
    redis : {
        host : process.env.REDIS_HOST,
        port: Number(process.env.REDIS_PORT),
        password: process.env.REDIS_PASSWORD,
        db: Number(process.env.REDIS_DB),
    },
    JWT_SECRET :{
        secret: process.env.JWT_SECRET
    },
    rabbitmq: {
        url: process.env.RABBITMQ_URL, 
        signupQueue: process.env.signupQueue,
        addElementQueue: process.env.addElementQueue,
        updateElementQueue: process.env.updateElementQueue,
        addAvatarQueue: process.env.addAvatarQueue,
        addMapQueue: process.env.addMapQueue
    },
    db : {
        user : process.env.user,
        host : process.env.host,
        database : process.env.database,
        password : process.env.password,
        port: Number(process.env.DB_PORT)
    },
    DB2 : {
        host : process.env.host,
        localDataCenter : process.env.localDataCenter,
    },
    AWS_S3:{
        accessKeyId: process.env.accessKeyId,
        secretAccessKey: process.env.secretAccessKey,
        region: process.env.region
    }
};
