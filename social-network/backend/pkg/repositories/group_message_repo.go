package repositories

import (
    "social-network/backend/pkg/models"
)

func SaveGroupMessage(groupID, senderID, content string) (*models.GroupMessage, error) {
    msg := &models.GroupMessage{
        ID:        generateUUID(),
        GroupID:   groupID,
        SenderID:  senderID,
        Content:   content,
        CreatedAt: now(),
    }
    _, err := DB.Exec(`
        INSERT INTO group_messages (id, group_id, sender_id, content, created_at)
        VALUES (?, ?, ?, ?, ?)`,
        msg.ID, msg.GroupID, msg.SenderID, msg.Content, msg.CreatedAt,
    )
    return msg, err
}

func GetGroupMessages(groupID string, limit, offset int) ([]models.GroupMessage, error) {
    rows, err := DB.Query(`
        SELECT id, group_id, sender_id, content, created_at
        FROM group_messages
        WHERE group_id = ?
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?`,
        groupID, limit, offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var msgs []models.GroupMessage
    for rows.Next() {
        var m models.GroupMessage
        if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
            return nil, err
        }
        msgs = append(msgs, m)
    }
    return msgs, nil
}

func IsUserInGroup(userID, groupID string) (bool, error) {
    var count int
    err := DB.QueryRow(`
        SELECT COUNT(*) FROM group_members
        WHERE group_id = ? AND user_id = ?`,
        groupID, userID,
    ).Scan(&count)
    return count > 0, err
}

func GetGroupMembers(groupID string) ([]string, error) {
    rows, err := DB.Query(`
        SELECT user_id FROM group_members
        WHERE group_id = ?`,
        groupID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var members []string
    for rows.Next() {
        var uid string
        if err := rows.Scan(&uid); err != nil {
            return nil, err
        }
        members = append(members, uid)
    }
    return members, nil
}