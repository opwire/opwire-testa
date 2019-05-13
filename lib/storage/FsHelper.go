package storage

var _current_, _history_ Fs

func SetFs(newFs Fs) {
	if newFs != nil {
		_history_ = _current_
		_current_ = newFs
	}
}

func GetFs() Fs {
	if _current_ == nil {
		_history_ = nil
		_current_ = NewOsFs()
	}
	return _current_
}

func Reset() {
	if _history_ != nil {
		_current_ = _history_
	}
}
