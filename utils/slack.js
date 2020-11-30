const Discord = require('discord.js');

const hook = new Discord.WebhookClient('780573485322862634', 'a1X_FBlH5D6VpRl311VPUoGeWnhinmCfg75Z2YFGLT7l9apSpjVSmU4gczcvf790j6O6');

hook.send('Test `message`!');

setTimeout(() => process.exit(1), 1000);