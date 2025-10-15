package httphandlers

type loginReq struct {
	Username string `json:"username" binding:"required,min=3,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type loginResp struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type itemReq struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"max=1000"`
	Quantity    int    `json:"quantity" binding:"required,min=0"`
}

type itemResp struct {
	ID          int    `json:"id"`
	Quantity    int    `json:"quantity"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type getItemHistoryResp struct {
	ItemID int               `json:"item_id"`
	Items  []itemHistoryResp `json:"items"`
}

type itemHistoryResp struct {
	UserID    int    `json:"user_id"`
	Operation string `json:"operation"`
	OldValue  string `json:"old_value,omitempty"`
	NewValue  string `json:"new_value,omitempty"`
	ChangedAt string `json:"changed_at"`
}
