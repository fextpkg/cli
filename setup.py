"""
Sometimes unconventional solutions are necessary to tackle unconventional situations.
We inherit from the wheel package to directly package the application into it,
avoiding the need to invent custom formats.
This allows for efficient utilization of a unified parser supported by Python itself.

Originally hosted on PyPI, but unfortunately removed for unknown reasons.
Thus, the name "fext-cli" is taken, ensuring no conflicts.
Currently, GitHub Releases are used to store the compiled application in a wheel package.

For initial installation, the script from "github.com/fextpkg/get" is used.
Self-update functionality is expected soon.
"""

import os
from dataclasses import dataclass

from setuptools import setup
from setuptools.command.install import install

# Execution file name.
EXE_NAME: str = "fext"

# Environment variables names.
# Package version.
ENV_VERSION: str = "FEXT_VERSION"
# Name of the directory with the executable files.
ENV_BINARY_DIR: str = "FEXT_BINARY_DIR"


def retrieve_env_variable(key: str) -> str:
    """
    Retrieves value from the environment variable.

    :raise RuntimeError: If it's not set.
    """
    if not (v := os.environ.get(key)):
        raise RuntimeError(f"Environment variable {key} is not specified")

    return v


@dataclass
class Platform:
    name: str
    executable_file: str

    @classmethod
    def windows(cls) -> "Platform":
        return cls(name="win_amd64", executable_file=f"{EXE_NAME}.exe")

    @classmethod
    def linux(cls) -> "Platform":
        return cls(name="manylinux_2_35_x86_64", executable_file=EXE_NAME)


class OverrideCommand(install):
    """
    Built-in setuptools commands don't support straightforward addition of binary files.
    More precisely, they **don't allow** adding them to scripts.
    We understand Python's stance on this matter,
    but we want to **avoid impacting** Python in any way because its execution consumes many resources.

    To address this issue, we modified this command to create an empty-package
    containing only metadata and a scripts directory with the binary file itself.

    Unfortunately, no builder can be configured as flexibly as setuptools itself.
    Consequently, none can support such commands without workarounds.
    It's not the best solution, but at least it's easy to maintain.

    Yes, direct invocation of ``setup.py`` is deprecated, but there's currently **no alternative**.
    """

    # Compiled binary file.
    executable_file: str
    # Working directory.
    source_dir: str = os.path.dirname(os.path.realpath(__file__))
    # Directory that contains compiled binary files.
    binary_dir: str = retrieve_env_variable(ENV_BINARY_DIR)

    def run(self):
        """
        The magical installation command that creates a bit of mess inside the package.
        """
        # As a precaution, run the original command just in case.
        super().run()

        # If the directory hasn't been created yet, create it.
        if not os.path.isdir(self.install_scripts):
            os.makedirs(self.install_scripts)

        # Specify both the external and internal paths to the executable file.
        source = os.path.join(self.source_dir, self.binary_dir, self.executable_file)
        target = os.path.join(self.install_scripts, self.executable_file)

        # If it happens that it already exists, remove it to avoid errors.
        if os.path.isfile(target):
            os.remove(target)

        # And perform a dirty trick.
        self.copy_file(source, target)


class Builder:
    def __init__(self) -> None:
        # Supported platforms for which the package needs to be built.
        self.platforms: list[Platform] = [
            Platform.windows(),
            Platform.linux(),
        ]

        # Prepare data
        self.version: str = retrieve_env_variable(ENV_VERSION)
        self.description, self.description_type = self.get_description()

    def _setup(self, platform_name: str) -> None:
        """
        Builds the package using ``setuptools``.
        """
        setup(
            # General information.
            name="fext-cli",
            version=self.version,
            description="Fast & Modern package manager",
            long_description=self.description,
            long_description_content_type=self.description_type,
            license="MIT",
            author="Andrew Krylov",
            author_email="any@lunte.dev",
            url="https://github.com/fextpkg/cli",
            keywords=["fast", "modern", "package", "manager"],
            # Ignore errors related to empty package
            # while simultaneously optimizing the package size.
            packages=[],
            # Leverage the ability to store external files.
            include_package_data=True,
            # Designate the package to be available only for a specific platform.
            options={"bdist_wheel": {"plat_name": platform_name}},
            # The magic lies in overriding the installation command.
            cmdclass={"install": OverrideCommand},
            # Don't generate in ".egg" format.
            zip_safe=False,
        )

    @staticmethod
    def get_description() -> tuple[str, str]:
        """
        Retrieves the text and type of README file.
        """
        with open("README.md", "r", encoding="utf-8") as f:
            return f.read(), "text/markdown"

    def pack(self, platform: Platform) -> None:
        """
        Builds the package for the specified platform.
        """
        # Adjust the path to the executable file according to the platform.
        OverrideCommand.executable_file = platform.executable_file
        self._setup(platform.name)

    def pack_all(self) -> None:
        """
        Builds the package for the all available platforms.
        """
        for plat in self.platforms:
            self.pack(plat)


if __name__ == "__main__":
    Builder().pack_all()
