# Digital Route Happy Push

An experiment

# Presentation

Hello and welcome to this presentation about the "HappyPush"

Happy push is a Feedback terminal inspired by the "HappyOrNot" terminal which can be found in alot of stores or servicecenters here in Sweden.

Our goal with HappyPush is to be able to measure Happines and the current mood in the office. As of today, we're using &frankly surveys every now and then asking the employees how they are doing. However, the answers on these surveys might be bery influenced by your current mood and not your "average" one. For example if I have a very good day and answers a &frankly survey, my answers will be more positive than if i have answerd they survey on a "bad day". 

With the HappyPush we're able to track happinies and the mood of the office in real time. This opens up for alot of possibilities:
With HappyPush, it's possible to correlate the office mood with certain events

This means its also possible to track the average happines and for example, send a notification to HR stating that the avrerage mood is below a certain threshold. This give HR the possibility to send out a &frankly survey to be able to find the root cause. 

HappyPush can also track how many people have answered "happy", how many have answerd "neutral" and how many that have answerd "sad" to be able to see if there are any big spans.

Technical stuff:

Three buttons: Green for happy/positive, Yellow for neutral and red for Sad/negative.
We're using ESP8266 connected to wifi. The ESP8266 is able to reconnect and alert via led if soething ist working. Like if it's disconnected from Wifi och can't push metrics.

The ESP talks to a API gateway over HTTPS using API-key as security. The API Gateway then triggers a lambda which currently stores the data in Cloudwatch. However, we can utilize several database backends like elasticsearch and dynamodb to be able to do even more complex queries.


HappyPush can be used to more than just measuring happinies in the office. As an example we can use the HappyPush to collect feedback about this FedEx. With that said, please rate your FedEx experience when you leave! The HappyPush is located right next to the door.
