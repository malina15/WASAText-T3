load("/scripts/lib/random.js");

// MongoDB Script: Populate Embedded Schema for WASAText (WORKING)
// Embedded model: conversations contain embedded messages (and embedded reactions/receipts).
// Run inside container mongosh: load("/scripts/populate_embedded.js")

const dbName = "wasatext_embedded";
use(dbName);

// Reset collections
db.users.drop();
db.conversations.drop();

print("🚀 Starting WASAText Embedded Schema Population");
print("📊 Generating synthetic messaging data...");

// ---------- helpers ----------

const firstNames = ["Ana","Maria","Ioana","Elena","Andrei","Mihai","Vlad","Radu","Irina","Teo","Daria","Alex","Matei","Sofia","Bianca","Paul"];
const lastNames  = ["Popescu","Ionescu","Georgescu","Dumitrescu","Marin","Stan","Matei","Ilie","Moldovan","Rusu","Toma","Petrescu"];

function makeUsername() {
  const f = choice(firstNames).toLowerCase();
  const l = choice(lastNames).toLowerCase();
  return `${f}.${l}${randInt(10,999)}`;
}

const USERS = 200;
const CONVERSATIONS = 80;

// ---------- 1) users ----------
let userDocs = [];
for (let i = 0; i < USERS; i++) {
  const username = makeUsername();
  userDocs.push({
    username,
    email: `${username}@example.com`,
    display_name: username,
    created_at: randomDateWithinDays(180),
    last_activity: randomDateWithinDays(14),
    is_online: Math.random() < 0.2
  });
}

db.users.insertMany(userDocs);
print(`✅ Inserted users: ${db.users.countDocuments()}`);

// cache user ids and usernames
const users = db.users.find({}, { _id: 1, username: 1, display_name: 1 }).toArray();
const userIds = users.map(u => u._id);

function pickParticipants() {
  const count = randInt(2, 6);
  const set = new Set();
  while (set.size < count) set.add(String(choice(userIds)));
  return Array.from(set).map(idStr => new ObjectId(idStr));
}

function userMini(uId) {
  const u = users.find(x => String(x._id) === String(uId));
  return { user_id: u._id, username: u.username, display_name: u.display_name };
}

const emojis = ["like", "love", "laugh", "wow", "sad", "angry"];
const sampleTexts = [
  "Hey! How are you?",
  "Ok, sounds good.",
  "I’ll send the details in a minute.",
  "Nice! 👍",
  "Let’s meet tomorrow.",
  "Can you review this?",
  "I agree with that.",
  "Great update, thanks!",
  "Any news on the task?",
  "Perfect, done."
];

// ---------- 2) conversations with embedded messages ----------
let convDocs = [];

for (let i = 0; i < CONVERSATIONS; i++) {
  const participants = pickParticipants(); // array of ObjectId
  const isGroup = participants.length > 2;

  const createdAt = randomDateWithinDays(120);
  const msgCount = randInt(50, 250); // enough messages to make counts meaningful
  let messages = [];

  // build messages chronologically
  let baseTime = new Date(createdAt.getTime());
  for (let m = 0; m < msgCount; m++) {
    // increase time a bit for each message
    baseTime = new Date(baseTime.getTime() + randInt(1, 300) * 1000); // +1..300 seconds

    const senderId = choice(participants);
    const msgId = new ObjectId();

    // reactions (0..3)
    let reactions = [];
    const rCount = randInt(0, 3);
    for (let r = 0; r < rCount; r++) {
      const reactorId = choice(participants);
      reactions.push({
        _id: new ObjectId(),
        user_id: reactorId,
        reaction_type: choice(emojis),
        timestamp: new Date(baseTime.getTime() + randInt(1, 60) * 1000)
      });
    }

    // receipts per participant (delivered/read)
    let receipts = [];
    for (const pid of participants) {
      const deliveredAt = new Date(baseTime.getTime() + randInt(1, 15) * 1000);
      const readChance = Math.random();
      receipts.push({
        user_id: pid,
        status: readChance < 0.75 ? "read" : "delivered",
        delivered_at: deliveredAt,
        read_at: readChance < 0.75 ? new Date(deliveredAt.getTime() + randInt(10, 600) * 1000) : null
      });
    }

    messages.push({
      _id: msgId,
      sender: userMini(senderId),
      content: choice(sampleTexts),
      message_type: "text",
      timestamp: new Date(baseTime),
      is_deleted: false,
      reactions,
      receipts
    });
  }

  const lastMessageAt = messages.length ? messages[messages.length - 1].timestamp : createdAt;

  convDocs.push({
    type: isGroup ? "group" : "direct",
    title: isGroup ? `Group Chat ${i + 1}` : null,
    participants: participants.map(pid => ({
      ...userMini(pid),
      joined_at: createdAt,
      role: "member"
    })),
    participants_count: participants.length,
    created_at: createdAt,
    last_message_at: lastMessageAt,
    message_count: msgCount,
    is_archived: false,
    messages
  });
}

db.conversations.insertMany(convDocs);

print(`✅ Inserted conversations: ${db.conversations.countDocuments()}`);
print("✅ Done. You can now run queries and explain() benchmarks.");
