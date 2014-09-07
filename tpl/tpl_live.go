// +build live

package tpl

const (
	debug = false
	template = `
<script>
const imgPath = 'image/'
const noTagDir = '%%NO_TAG_DIR%%'

var preloadImg = $('<img>')

var images = []

//            Q   A   W   S   E   D   R   F   Z   X   C   V   1   2   3   4
const keys = [81, 65, 87, 83, 69, 68, 82, 70, 90, 88, 67, 86, 49, 50, 51, 52]
var tags = []

var sorted = []
var tagsCounts = []
var randomQueue = []
var currentImg = 0

$(function() {
	$(window).keydown(function(e) {
		var i = -1
		// shift + space
		if (e.keyCode == 32 && e.shiftKey) {
			show(SHOW.PREV)
		// ctrl + W
		} else if (e.keyCode == 87 && e.ctrlKey) {
			$.get('/close')
			return true
		// Enter
		} else if (e.keyCode == 13) {
			var mode = BUILD.VIEWED
			if (e.ctrlKey) {
				mode = BUILD.ALL
			} else if (e.shiftKey) {
				mode = BUILD.TAGGED
			}

			$.post('/save', buildResult(mode), function() {
				loadImages()
			})
		} else if (e.altKey || e.ctrlKey || e.shiftKey) {
			return true
		//         space              pgdown             right
		} else if (e.keyCode == 32 || e.keyCode == 34 || e.keyCode == 39) {
			show(SHOW.NEXT)
		//         backspace         pgup               left
		} else if (e.keyCode == 8 || e.keyCode == 33 || e.keyCode == 37) {
			show(SHOW.PREV)
		// home
		} else if (e.keyCode == 36) {
			show(SHOW.FIRST)
		// end
		} else if (e.keyCode == 35) {
			show(SHOW.LAST)
		// /
		} else if (e.keyCode == 111) {
			show(SHOW.RANDOM)
		// esc
		} else if (e.keyCode == 27) {
			$.get('/close')
		// numpad +
		} else if (e.keyCode == 107 && tags.length < keys.length) {
			var tag = prompt('New tag:')

			if (typeof tag != 'undefined' && tag != null && tag.length > 0) {
				tags.push(tag)
				drawKeymap()
			}
		// ~
		} else if (e.keyCode == 192) {
			$('.overlay').fadeToggle(150)
		// ins
		} else if (e.keyCode == 45) {
			if (typeof sorted[currentImg] != 'undefined') {
				tagsCounts[sorted[currentImg]]--
			}

			if (typeof sorted[currentImg] == 'undefined' || sorted[currentImg] >= tags.length - 1) {
				sorted[currentImg] = 0
			} else {
				sorted[currentImg]++
			}
			tagsCounts[sorted[currentImg]]++

			show()
			drawKeymap()
		// del
		} else if (e.keyCode == 46) {
			tagsCounts[sorted[currentImg]]--
			delete sorted[currentImg]
			show()
			drawKeymap()
		// *
		} else if (e.keyCode == 106) {
			randomize()
			reset()
			show(SHOW.FIRST)
			drawKeymap()
		} else if ((i = keys.indexOf(e.keyCode)) >= 0 && typeof tags[i] != 'undefined') {
			if (sorted[currentImg] != i) {
				if (typeof sorted[currentImg] != 'undefined') {
					tagsCounts[sorted[currentImg]]--
				}

				tagsCounts[i]++
				sorted[currentImg] = i
			} else {
				tagsCounts[i]--
				delete sorted[currentImg]
			}

			show()
			drawKeymap()
		} else {
			console.log(e.keyCode)
			return true
		}

		return false
	})
	.bind('mousewheel', function(e) {
		if (e.originalEvent.detail >= 0) {
			show(SHOW.NEXT)
		} else {
			show(SHOW.PREV)
		}
	})
	.resize(resize).triggerHandler('resize')

	$('#current_image').load(calcImgSize)

	loadImages()
})

function resize() {
	$('#current_image').css({
		maxHeight: window.innerHeight - 4,
		maxWidth: window.innerWidth - 4
	})
	$('#current_flash embed').css({
		height: window.innerHeight - 4,
		width: window.innerWidth - 4
	})
	calcImgSize()
}

function calcImgSize() {
	if (/\.swf$/.test(images[currentImg])) {
		$('#size').text('')
	} else {
		var img = $('#current_image')[0]
		var natural = img.naturalWidth * img.naturalHeight
		var actual = img.width * img.height
		var scale = actual / natural * 100

		$('#size').text(img.naturalWidth + '×' + img.naturalHeight + ' ' + scale.toFixed(0) + '%')
	}
	$('#loading').hide()
}

function loadImages() {
	$.get('/images', function(data) {
		images = data

		reset()
		show(SHOW.FIRST)

		if (images.length == 0) {
			$.get('/close')
			return
		}

		loadTags()
	})
}

function loadTags() {
	$.get('/tags', function(data) {
		tags = data

		drawKeymap()
	})
}

function reset() {
	sorted = []
	randomQueue = []

	tagsCounts = []
	for (var i in keys) {
		tagsCounts.push(0)
	}
}

function drawKeymap() {
	$('#keymap').empty()

	var totalCount = 0
	for (var i in keys) {
		if (typeof tags[i] == 'undefined') {
			break
		}

		var count = ''
		if (tagsCounts[i] > 0) {
			totalCount += tagsCounts[i]
			count = ' (' + tagsCounts[i] + ')'
		}

		var tagText = String.fromCharCode(keys[i]) + ': ' + tags[i] + count
		var tag = $('<div>')
			.text(tagText)
			.appendTo('#keymap')

		if (sorted[currentImg] == i) {
			tag.addClass('current_tag')
		}
	}

	$('<div class=total_tags>')
		.text('total: ' + totalCount)
		.appendTo('#keymap')
	$('<div>')
		.text('skipped: ' + (sorted.length - totalCount))
		.appendTo('#keymap')
}

const SHOW = {
	CURRENT: 0,
	FIRST:   1,
	PREV:    2,
	NEXT:    3,
	LAST:    4,
	RANDOM:  5
}

function show(direction) {
	if (images.length == 0) {
		$('body').text('No images here')
		return
	}

	switch (direction) {
	case SHOW.FIRST:
		currentImg = 0
		break
	case SHOW.PREV:
		if (currentImg == 0) {
			return
		}

		currentImg -= 1
		break
	case SHOW.NEXT:
		if (currentImg == images.length - 1) {
			return
		}

		currentImg += 1
		break
	case SHOW.LAST:
		currentImg = images.length - 1
		break
	case SHOW.RANDOM:
		currentImg = null
		while (currentImg == null || randomQueue.indexOf(currentImg) >= 0) {
			currentImg = Math.floor(Math.random() * images.length)
		}

		break
	}

	if (currentImg < 0) {
		currentImg = images.length - 1
	} else if (currentImg >= images.length) {
		currentImg = 0
	}

	randomQueue.push(currentImg)

	var diff = Math.round(randomQueue.length - (images.length / 2))
	if (diff > 0) {
		randomQueue.splice(0, diff)
	}

	var name = images[currentImg]
	if (/\.swf$/.test(name)) {
		$('#loading').hide()
		$('#current_flash').show()
		$('#current_flash embed').attr('src', imgPath + name + '?_=' + Math.random())
		$('#current_image').hide()
		$('#size').text('')
	} else {
		$('#loading').show()
		$('#current_image').attr('src', imgPath + name).show()
		$('#current_flash').hide()

		if (currentImg < images.length - 1) {
			var preloadName = images[currentImg + 1]
			preloadImg.attr('src', imgPath + preloadName)
		}
	}

	$('#name').text(name)
	$('#counts').text((currentImg + 1) + '/' + images.length)
	drawKeymap()
}

const BUILD = {
	ALL:    0,
	TAGGED: 1,
	VIEWED: 2,
}

function buildResult(mode) {
	var result = {}

	var count = sorted.length
	if (mode == BUILD.ALL)
		count = images.length

	for (var i = 0; i < count; i++) {
		if (typeof sorted[i] != 'undefined') {
			result[images[i]] = tags[sorted[i]]
		} else if (mode == BUILD.VIEWED || mode == BUILD.ALL) {
			result[images[i]] = noTagDir
		}
	}

	return result
}

function randomize() {
	var sorted = []

	while (images.length) {
		var buf = images.splice(Math.random() * images.length, 1)
		sorted.push(buf)
	}

	images = sorted
}
</script>
<style>
body {
	font-family: Verdana, Helvetica, Arial;
	margin: 2px;
	color: #333;
}

#name {
	font-size: 14px;
	margin-right: 10px;
}

#counts {
	margin-right: 10px;
}

#keymap {
	display: inline-block;
}

.overlay {
	position: absolute;
	margin: 3px;
}

.overlay .panel {
	background-color: #F7F7F7;
	padding: 3px;
	font-size: 12px;
	border: 1px solid;
	border-radius: 3px;
	margin-bottom: 5px;
}

.current_tag {
	color: #333;
	font-weight: bold;
}

.total_tags {
	margin-top: 5px;
}

#current_image, #current_flash {
	margin: 0 auto;
	display: block;
}
</style>

<div class=overlay>
	<div class=panel>
		<span id=name></span>
		<span id=counts></span>
		<span id=size></span>
		<span id=loading>♥</span>
	</div>
	<div id=keymap class=panel></div>
</div>

<img id=current_image>
<object id=current_flash><embed></object>
`
)