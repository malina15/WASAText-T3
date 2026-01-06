# Analiză Detaliată a Discrepanțelor - Proiect WASA

## Rezumat Executiv

Proiectul actual este un **WASAphoto** (aplicație de partajare foto), dar cerințele specifice sunt pentru **WASAText** (aplicație de mesagerie). Există multiple discrepanțe critice între implementarea actuală și cerințele din specificația WASAText.

---

## 1. DISCREPANȚE CRITICE - OpenAPI Specification

### 1.1 Schema de Login Incorectă

**Cerință WASAText:**
```yaml
requestBody:
  content:
    application/json:
      schema:
        type: object
        properties:
          name:  # ← TREBUIE să fie "name", nu "identifier"
            type: string
            example: Maria
            pattern: '^.*?$'
            minLength: 3
            maxLength: 16
```

**Implementare Actuală (doc/api.yaml:801-806):**
```yaml
user_login:
  properties:
    identifier:  # ← GREȘIT! Trebuie "name"
      $ref: "#/components/schemas/user_nickname/properties/nickname"
```

**Backend Code (service/api/session.go:16-17):**
```go
var user User
err := json.NewDecoder(r.Body).Decode(&user)
// User struct expects "user_id" field, not "name"
```

**Problema:** 
- OpenAPI folosește `identifier` în loc de `name`
- Backend-ul așteaptă `user_id` în JSON
- Specificația WASAText cere explicit `name`

---

### 1.2 OperationIds Lipsă din OpenAPI (dar implementate în backend)

#### 1.2.1 `getMyConversations`
- **Status:** ❌ Lipsește din OpenAPI
- **Backend:** ✅ Implementat ca `/users/:id/chats` (service/api/api-handler.go:32)
- **Handler:** `listChats` (service/api/chat_handlers.go:13)
- **Cerință WASAText:** Trebuie să returneze lista de conversații sortate invers cronologic

#### 1.2.2 `getConversation`
- **Status:** ❌ Lipsește din OpenAPI
- **Backend:** ✅ Implementat ca `/users/:id/chats/:peer/messages` (service/api/api-handler.go:33)
- **Handler:** `listMessages` (service/api/chat_handlers.go:31)
- **Cerință WASAText:** Trebuie să returneze mesajele dintr-o conversație

#### 1.2.3 `sendMessage`
- **Status:** ❌ Lipsește din OpenAPI
- **Backend:** ✅ Implementat ca `POST /users/:id/chats/:peer/messages` (service/api/api-handler.go:34)
- **Handler:** `sendMessage` (service/api/chat_handlers.go:50)
- **Cerință WASAText:** Trebuie să permită trimiterea de mesaje text sau GIF

---

### 1.3 OperationIds Complet Lipsă (nici implementate, nici în OpenAPI)

#### 1.3.1 `forwardMessage`
- **Status:** ❌ Complet lipsă
- **Cerință WASAText:** Utilizatorul trebuie să poată forwarda mesaje

#### 1.3.2 `commentMessage`
- **Status:** ❌ Complet lipsă
- **Notă:** Există `commentPhoto` pentru comentarii pe poze, dar NU pentru reacții (emoticoane) pe mesaje
- **Cerință WASAText:** Utilizatorii trebuie să poată reacționa la mesaje cu emoticoane

#### 1.3.3 `uncommentMessage`
- **Status:** ❌ Complet lipsă
- **Notă:** Există `uncommentPhoto` pentru poze, dar NU pentru mesaje
- **Cerință WASAText:** Utilizatorii trebuie să poată șterge propriile reacții

#### 1.3.4 `deleteMessage`
- **Status:** ❌ Complet lipsă
- **Cerință WASAText:** Utilizatorii trebuie să poată șterge mesajele trimise

#### 1.3.5 `addToGroup`
- **Status:** ❌ Complet lipsă
- **Cerință WASAText:** Membrii grupului trebuie să poată adăuga alți utilizatori

#### 1.3.6 `leaveGroup`
- **Status:** ❌ Complet lipsă
- **Cerință WASAText:** Utilizatorii trebuie să poată părăsi grupuri

#### 1.3.7 `setGroupName`
- **Status:** ❌ Complet lipsă
- **Cerință WASAText:** Trebuie să existe posibilitatea de a seta/modifica numele grupului

#### 1.3.8 `setMyPhoto`
- **Status:** ❌ Complet lipsă
- **Notă:** Există `uploadPhoto` pentru poze de profil, dar nu există un endpoint dedicat pentru foto de profil
- **Cerință WASAText:** Utilizatorii trebuie să poată seta foto de profil

#### 1.3.9 `setGroupPhoto`
- **Status:** ❌ Complet lipsă
- **Cerință WASAText:** Trebuie să existe posibilitatea de a seta foto pentru grupuri

---

## 2. DISCREPANȚE CRITICE - Funcționalități WASAText

### 2.1 Grupuri (Groups)

**Status:** ❌ Complet lipsă

**Ce lipsește:**
- Tabele de bază de date pentru grupuri
- API endpoints pentru grupuri
- Funcționalități: creare grup, adăugare membri, părăsire grup, setare nume/foto grup
- Frontend pentru gestionarea grupurilor

**Cerință WASAText:**
> "The user can create a new group with any number of other WASAText users to start a conversation. Group members can add other users to the group, but users cannot join groups on their own or even see groups they aren't a part of. Additionally, users have the option to leave a group at any time."

---

### 2.2 Reacții la Mesaje (Message Reactions/Comments)

**Status:** ❌ Complet lipsă

**Ce lipsește:**
- Tabele de bază de date pentru reacții la mesaje
- API endpoints pentru `commentMessage` și `uncommentMessage`
- Frontend pentru afișarea și adăugarea reacțiilor

**Cerință WASAText:**
> "Users can also react to messages (a.k.a. comment them) with an emoticon, and delete their reactions at any time (a.k.a. uncomment)."

**Notă:** Există comentarii pe poze (`commentPhoto`), dar acestea sunt diferite de reacțiile pe mesaje.

---

### 2.3 Forward Mesaje

**Status:** ❌ Complet lipsă

**Cerință WASAText:**
> "The user can send a new message, reply to an existing one, forward a message, and delete any sent messages."

---

### 2.4 Reply la Mesaje

**Status:** ❌ Complet lipsă

**Cerință WASAText:**
> "The user can send a new message, reply to an existing one, forward a message, and delete any sent messages."

---

### 2.5 Delete Mesaje

**Status:** ❌ Complet lipsă

**Cerință WASAText:**
> "The user can send a new message, reply to an existing one, forward a message, and delete any sent messages."

---

### 2.6 Read Receipts (Checkmarks)

**Status:** ❌ Complet lipsă

**Cerință WASAText:**
> "Each message includes the timestamp, the content (whether text or photo), and the sender's username for received messages, or one/two checkmarks to indicate the status of sent messages. One checkmark indicates that the message has been received by the recipient (by all the recipients for groups) in their conversation list. Two checkmarks mean that the message has been read by the recipient (by all the recipients for groups) within the conversation itself."

**Ce lipsește:**
- Câmpuri în baza de date pentru status mesaj (sent, received, read)
- Logică pentru tracking read receipts
- Frontend pentru afișarea checkmarks

---

### 2.7 Preview Mesaje în Lista de Conversații

**Status:** ⚠️ Parțial implementat

**Cerință WASAText:**
> "Each element in the list must display the username of the other person or the group name, the user profile photo or the group photo, the date and time of the latest message, the preview (snippet) of the text message, or an icon for a photo message."

**Implementare Actuală (webui/src/views/ChatsView.vue):**
- ✅ Afișează username
- ❌ NU afișează foto de profil
- ❌ NU afișează data/ora ultimului mesaj
- ❌ NU afișează preview/snippet al mesajului
- ❌ NU suportă grupuri

---

### 2.8 Sortare Conversații

**Status:** ❌ Nu este implementată sortarea invers cronologică

**Cerință WASAText:**
> "The user is presented with a list of conversations with other users or with groups, sorted in reverse chronological order."

**Implementare Actuală (service/api/chat_handlers.go:20):**
```go
following, err := rt.db.GetFollowing(database.User{IdUser: requester})
// Returnează doar lista de utilizatori urmăriți, NU conversațiile sortate
```

---

### 2.9 Foto de Profil Utilizator

**Status:** ❌ Complet lipsă

**Cerință WASAText:**
- Utilizatorii trebuie să poată seta foto de profil (`setMyPhoto`)
- Foto de profil trebuie afișată în lista de conversații

---

### 2.10 Foto de Profil Grup

**Status:** ❌ Complet lipsă

**Cerință WASAText:**
- Grupurile trebuie să aibă foto (`setGroupPhoto`)
- Foto grup trebuie afișată în lista de conversații

---

## 3. DISCREPANȚE TEHNICE

### 3.1 CORS Max-Age

**Status:** ❌ Nu este setat la 1 secundă în cod

**Cerință WASAText:**
> "To avoid problems during the homework grading, you should allow all origins and you should set the 'Max-Age' attribute to 1 second."

**Implementare Actuală (cmd/webapi/cors.go:12-22):**
```go
func applyCORSHandler(h http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{...}),
		handlers.AllowedMethods([]string{...}),
		handlers.AllowedOrigins([]string{"*"}),  // ✅ OK
		// ❌ LIPSEȘTE: handlers.MaxAge(1)
	)(h)
}
```

**Soluție:** Trebuie adăugat `handlers.MaxAge(1)` în configurația CORS.

---

### 3.2 Linting Configuration Files

**Status:** ❌ Toate fișierele de configurare pentru linting lipsesc

**Fișiere Lipsă:**
1. `.spectral.js` - pentru validarea OpenAPI
2. `.golangci.yml` - pentru linting Go
3. `webui/eslint.config.mjs` - pentru linting Vue/JavaScript

**Cerință:** Aceste fișiere sunt necesare pentru validarea codului conform standardelor cerute.

---

## 4. DISCREPANȚE ÎN STRUCTURA BAZEI DE DATE

### 4.1 Tabele Lipsă pentru WASAText

**Tabele necesare care lipsesc:**
1. `groups` - pentru stocarea grupurilor
2. `group_members` - pentru relația many-to-many între utilizatori și grupuri
3. `message_reactions` - pentru reacțiile (emoticoane) pe mesaje
4. `message_status` - pentru read receipts (sent, received, read)
5. `user_profile_photos` - pentru foto de profil utilizator
6. `group_photos` - pentru foto de grup

**Tabele existente care sunt OK:**
- ✅ `users` - pentru utilizatori
- ✅ `messages` - pentru mesaje directe (dar fără status/reactions)

---

## 5. DISCREPANȚE ÎN FRONTEND

### 5.1 Funcționalități Lipsă în Frontend

1. **Grupuri:**
   - ❌ Creare grup
   - ❌ Adăugare membri la grup
   - ❌ Părăsire grup
   - ❌ Setare nume/foto grup
   - ❌ Afișare grupuri în lista de conversații

2. **Mesaje:**
   - ❌ Forward mesaje
   - ❌ Reply la mesaje
   - ❌ Delete mesaje
   - ❌ Reacții (emoticoane) pe mesaje
   - ❌ Afișare read receipts (checkmarks)
   - ❌ Preview mesaje în lista de conversații
   - ❌ Sortare conversații invers cronologic

3. **Profil:**
   - ❌ Setare foto de profil
   - ❌ Afișare foto de profil în conversații

---

## 6. REZUMAT OPERATIONIDs

### 6.1 OperationIds Existente și Corecte
- ✅ `doLogin` - există (dar schema este greșită)
- ✅ `setMyUserName` - există

### 6.2 OperationIds Implementate dar Lipsă din OpenAPI
- ⚠️ `getMyConversations` - backend OK, lipsește din OpenAPI
- ⚠️ `getConversation` - backend OK, lipsește din OpenAPI
- ⚠️ `sendMessage` - backend OK, lipsește din OpenAPI

### 6.3 OperationIds Complet Lipsă
- ❌ `forwardMessage`
- ❌ `commentMessage`
- ❌ `uncommentMessage`
- ❌ `deleteMessage`
- ❌ `addToGroup`
- ❌ `leaveGroup`
- ❌ `setGroupName`
- ❌ `setMyPhoto`
- ❌ `setGroupPhoto`

---

## 7. PRIORITIZARE REPARĂRI

### Prioritate CRITICĂ (Blocante pentru evaluare):
1. ✅ Corectare schema login (`name` în loc de `identifier`)
2. ✅ Adăugare CORS Max-Age = 1 secundă
3. ✅ Adăugare operationIds lipsă în OpenAPI pentru endpoint-urile existente
4. ✅ Creare fișiere de configurare linting

### Prioritate ÎNALTĂ (Funcționalități WASAText esențiale):
5. ✅ Implementare grupuri (bază de date + API + frontend)
6. ✅ Implementare `setMyPhoto` și `setGroupPhoto`
7. ✅ Implementare `forwardMessage`, `deleteMessage`, reply
8. ✅ Implementare `commentMessage` și `uncommentMessage` (reacții)
9. ✅ Implementare read receipts (checkmarks)

### Prioritate MEDIE (Îmbunătățiri UX):
10. ✅ Sortare conversații invers cronologic
11. ✅ Preview mesaje în lista de conversații
12. ✅ Afișare foto profil în conversații

---

## 8. NOTĂ IMPORTANTĂ

**Proiectul actual este WASAphoto (photo sharing), dar cerințele sunt pentru WASAText (messaging).** 

Aceasta înseamnă că:
- Multe funcționalități existente (poze, like-uri, comentarii pe poze) NU sunt relevante pentru WASAText
- Trebuie implementate funcționalități complet noi (grupuri, reacții pe mesaje, read receipts)
- Structura bazei de date trebuie extinsă semnificativ
- Frontend-ul trebuie reproiectat pentru a se concentra pe mesagerie, nu pe partajare foto

---

## 9. VERIFICARE FINALĂ

După implementarea tuturor reparațiilor, verificați:

- [ ] Toate operationIds cerute există în OpenAPI
- [ ] Schema login folosește `name` (nu `identifier`)
- [ ] CORS Max-Age = 1 secundă
- [ ] Fișiere linting configurate și funcționale
- [ ] Grupurile funcționează complet
- [ ] Reacțiile pe mesaje funcționează
- [ ] Read receipts funcționează
- [ ] Forward/Reply/Delete mesaje funcționează
- [ ] Foto profil utilizator și grup funcționează
- [ ] Frontend afișează toate funcționalitățile WASAText

---

**Data analizei:** 2025-01-XX
**Versiune specificație WASAText:** Version 1

