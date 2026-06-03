# CLAUDE.md

## P0 Optimization Guidelines

- Question detail image IDs must be built without placeholder zeros; initialize ID slices with capacity when appending.
- List endpoints should avoid per-question file lookups. Use `service.QuestionUtil()` batch helpers to resolve question images and answer avatars.
- Answer avatar lists should never fetch or return more than `consts.MaxAvatarsPerQuestion` per question unless the API explicitly requires all avatars.
- Notification queries should not reuse a chained model across different notification types. Build independent queries or aggregate by type.
- Notification rendering must tolerate missing related questions, answers, or users instead of panicking on map dereferences.
- Keep generated DAO/DO/entity files untouched; add reusable business helpers in existing service/logic layers.
