document.addEventListener('DOMContentLoaded', function () {
	document.getElementById('uploadForm').addEventListener('submit', function (e) {
		e.preventDefault();
		const currentPath = encodeURIComponent(window.location.pathname);
		this.action = `/upload?path=${currentPath}`;
		this.submit();
	});

	document.getElementById('createDir').addEventListener('submit', function (e) {
		e.preventDefault();
		const currentPath = encodeURIComponent(window.location);
		this.action = `/create_dir?path=${currentPath}`;
		this.submit();
	});
});
