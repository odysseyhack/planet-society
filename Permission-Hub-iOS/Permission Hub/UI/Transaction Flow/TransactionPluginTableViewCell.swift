//
//  TransactionPluginTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionPluginTableViewCell: PHBaseTableViewCell {

    // MARK: - Private properties

    private lazy var stackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.spacing = 20
        stackView.alignment = .center

        return stackView
    }()

    private lazy var pluginButton: UIButton = {

        let button = UIButton()
        button.translatesAutoresizingMaskIntoConstraints = false

        button.setContentHuggingPriority(.required, for: .horizontal)

        button.addTarget(
            self,
            action: #selector(pluginButtonTapped),
            for: .touchUpInside)

        return button
    }()

    private lazy var descriptionLabel: UILabel = {

        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false
        label.font = PHFonts.regular()
        label.textColor = PHColors.greyishBrown
        label.numberOfLines = 0

        return label
    }()

    private var callback: (() -> Void)?

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    private func configure() {

        selectionStyle = .none

        addSubview(stackView)

        let margin: CGFloat = 15
        stackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        stackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        stackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true
        stackView.bottomAnchor.constraint(equalTo: bottomAnchor, constant: -margin).isActive = true

        stackView.addArrangedSubview(pluginButton)
        stackView.addArrangedSubview(descriptionLabel)
    }

    func configure(withImage
        image: UIImage?, andText
        text: String,
        callback: @escaping () -> Void) {

        pluginButton.setImage(image, for: .normal)
        descriptionLabel.text = text
        self.callback = callback
        
        setNeedsLayout()
    }

    // MARK: - Selectors

    @objc private func pluginButtonTapped(_ sender: UIButton) {
        callback?()
    }
}
