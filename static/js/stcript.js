document.getElementById('uploadForm').addEventListener('submit', function () {
	const currentPath = encodeURIComponent(window.location.pathname);
	this.action = `/upload?path=${currentPath}`;
});

document.getElementById('createDir').addEventListener('submit', function () {
	const currentPath = encodeURIComponent(window.location.pathname);
	this.action = `/create_dir?path=${currentPath}`;
});
