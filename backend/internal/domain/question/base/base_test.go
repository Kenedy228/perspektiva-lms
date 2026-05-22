package base

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	media2 "gitflic.ru/lms/backend/internal/domain/shared/media"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

func mustTitle(t *testing.T, v string) title.Title {
	t.Helper()

	got, err := title.New(v)
	if err != nil {
		t.Fatalf("title.New() error = %v", err)
	}

	return got
}

func mustAttachment(t *testing.T, mt media2.Type, key string, size int64) attachment.Attachment {
	t.Helper()

	f, err := file.New(key, size)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	m, err := media2.New(mt, f)
	if err != nil {
		t.Fatalf("media.New() error = %v", err)
	}

	a, err := attachment.New(m)
	if err != nil {
		t.Fatalf("attachment.New() error = %v", err)
	}

	return a
}

func TestBase_Attachment(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	tests := []struct {
		name   string
		fields fields
		want   attachment.Attachment
		want1  bool
	}{
		{
			name: "returns attachment and true when attachment exists",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
				hasAttachment: true,
			},
			want:  mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
			want1: true,
		},
		{
			name: "returns zero attachment and false when attachment does not exist",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
			want:  attachment.Attachment{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			got, got1 := b.Attachment()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Attachment() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Attachment() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBase_ChangeAttachment(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	type args struct {
		a attachment.Attachment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "sets attachment for base without attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
			args: args{
				a: mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
			},
		},
		{
			name: "replaces existing attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    mustAttachment(t, media2.TypeAudio, "content/questions/audio/task.mp3", 2048),
				hasAttachment: true,
			},
			args: args{
				a: mustAttachment(t, media2.TypeImage, "content/questions/image/updated.jpg", 4096),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			b.ChangeAttachment(tt.args.a)

			if !reflect.DeepEqual(b.attachment, tt.args.a) {
				t.Errorf("ChangeAttachment() attachment = %v, want %v", b.attachment, tt.args.a)
			}
			if !b.hasAttachment {
				t.Errorf("ChangeAttachment() hasAttachment = %v, want %v", b.hasAttachment, true)
			}
		})
	}
}

func TestBase_ChangeTitle(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	type args struct {
		t title.Title
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "changes title",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Старый заголовок"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
			args: args{
				t: mustTitle(t, "Новый заголовок"),
			},
		},
		{
			name: "changes title and keeps attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос до изменения"),
				attachment:    mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
				hasAttachment: true,
			},
			args: args{
				t: mustTitle(t, "Вопрос после изменения"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			oldAttachment := b.attachment
			oldHasAttachment := b.hasAttachment

			b.ChangeTitle(tt.args.t)

			if !reflect.DeepEqual(b.title, tt.args.t) {
				t.Errorf("ChangeTitle() title = %v, want %v", b.title, tt.args.t)
			}
			if !reflect.DeepEqual(b.attachment, oldAttachment) {
				t.Errorf("ChangeTitle() attachment changed = %v, want %v", b.attachment, oldAttachment)
			}
			if b.hasAttachment != oldHasAttachment {
				t.Errorf("ChangeTitle() hasAttachment changed = %v, want %v", b.hasAttachment, oldHasAttachment)
			}
		})
	}
}

func TestBase_Clone(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *Base
	}{
		{
			name: "clone without attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
			want: func() *Base {
				id := uuid.New()
				ttl := mustTitle(t, "placeholder")
				_ = id
				_ = ttl
				return nil
			}(),
		},
		{
			name: "clone with attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    mustAttachment(t, media2.TypeAudio, "content/questions/audio/task.mp3", 2048),
				hasAttachment: true,
			},
			want: func() *Base {
				id := uuid.New()
				ttl := mustTitle(t, "placeholder")
				_ = id
				_ = ttl
				return nil
			}(),
		},
	}
	for i := range tests {
		tests[i].want = &Base{
			id:            tests[i].fields.id,
			title:         tests[i].fields.title,
			attachment:    tests[i].fields.attachment,
			hasAttachment: tests[i].fields.hasAttachment,
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			if got := b.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
			if got := b.Clone(); got == b {
				t.Errorf("Clone() returned same pointer")
			}
		})
	}
}

func TestBase_HasAttachment(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returns true when attachment exists",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
				hasAttachment: true,
			},
			want: true,
		},
		{
			name: "returns false when attachment does not exist",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			if got := b.HasAttachment(); got != tt.want {
				t.Errorf("HasAttachment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase_ID(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns id",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
		},
		{
			name: "returns id with attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    mustAttachment(t, media2.TypeAudio, "content/questions/audio/task.mp3", 2048),
				hasAttachment: true,
			},
		},
	}
	for i := range tests {
		tests[i].want = tests[i].fields.id
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			if got := b.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase_RemoveAttachment(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "removes existing attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
				hasAttachment: true,
			},
		},
		{
			name: "remove when attachment already absent",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			b.RemoveAttachment()

			if !reflect.DeepEqual(b.attachment, attachment.Attachment{}) {
				t.Errorf("RemoveAttachment() attachment = %v, want %v", b.attachment, attachment.Attachment{})
			}
			if b.hasAttachment {
				t.Errorf("RemoveAttachment() hasAttachment = %v, want %v", b.hasAttachment, false)
			}
		})
	}
}

func TestBase_Title(t *testing.T) {
	type fields struct {
		id            uuid.UUID
		title         title.Title
		attachment    attachment.Attachment
		hasAttachment bool
	}
	tests := []struct {
		name   string
		fields fields
		want   title.Title
	}{
		{
			name: "returns title",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 1"),
				attachment:    attachment.Attachment{},
				hasAttachment: false,
			},
			want: mustTitle(t, "Вопрос 1"),
		},
		{
			name: "returns title with attachment",
			fields: fields{
				id:            uuid.New(),
				title:         mustTitle(t, "Вопрос 2"),
				attachment:    mustAttachment(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
				hasAttachment: true,
			},
			want: mustTitle(t, "Вопрос 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:            tt.fields.id,
				title:         tt.fields.title,
				attachment:    tt.fields.attachment,
				hasAttachment: tt.fields.hasAttachment,
			}
			if got := b.Title(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		t title.Title
	}
	tests := []struct {
		name    string
		args    args
		want    *Base
		wantErr bool
	}{
		{
			name: "creates base with title and generated id",
			args: args{
				t: mustTitle(t, "Новый вопрос"),
			},
			wantErr: false,
		},
		{
			name: "creates base with another title and empty attachment",
			args: args{
				t: mustTitle(t, "Вопрос с пустым вложением"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got == nil {
					t.Fatalf("New() got nil, want non-nil")
				}
				if got.id == uuid.Nil {
					t.Errorf("New() id = %v, want non-nil uuid", got.id)
				}
				if !reflect.DeepEqual(got.title, tt.args.t) {
					t.Errorf("New() title = %v, want %v", got.title, tt.args.t)
				}
				if !reflect.DeepEqual(got.attachment, attachment.Attachment{}) {
					t.Errorf("New() attachment = %v, want %v", got.attachment, attachment.Attachment{})
				}
				if got.hasAttachment {
					t.Errorf("New() hasAttachment = %v, want %v", got.hasAttachment, false)
				}
			}
		})
	}
}
